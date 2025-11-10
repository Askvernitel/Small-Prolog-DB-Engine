:- use_module(library(http/thread_httpd)).
:- use_module(library(http/http_dispatch)).
:- use_module(library(http/http_parameters)).
:- use_module(library(http/http_json)).
:- dynamic table_schema/2.
:- dynamic table_data/3.

db_directory('db_files/').
server_port(8081).

init_db :-
    db_directory(Dir),
    (exists_directory(Dir) -> true ; make_directory(Dir)).

:- http_handler(root(query), handle_query, []).

start_server :-
    init_db,
    server_port(Port),
    http_server(http_dispatch, [port(Port)]),
    format('Database server running on http://localhost:~w~n', [Port]),
    format('Send POST requests to /query with JSON body~n', []).

stop_server :-
    server_port(Port),
    http_stop_server(Port, []).

handle_query(Request) :-
    http_read_json_dict(Request, QueryDict),
    process_query(QueryDict, Response),
    reply_json_dict(Response).

process_query(Dict, Response) :-
    Type = Dict.get(type),
    (   Type = "create_table"
    ->  create_table_handler(Dict, Response)
    ;   Type = "insert"
    ->  insert_handler(Dict, Response)
    ;   Type = "select"
    ->  select_handler(Dict, Response)
    ;   Type = "update"
    ->  update_handler(Dict, Response)
    ;   Type = "delete"
    ->  delete_handler(Dict, Response)
    ;   Response = _{status: "error", message: "Unknown query type"}
    ).

create_table_handler(Dict, Response) :-
    Table = Dict.get(table),
    Columns = Dict.get(columns),
    (   table_schema(Table, _)
    ->  Response = _{status: "error", message: "Table already exists"}
    ;   assert(table_schema(Table, Columns)),
        save_schema(Table),
        Response = _{status: "success", message: "Table created", table: Table}
    ).

insert_handler(Dict, Response) :-
    Table = Dict.get(table),
    Values = Dict.get(values),
    (   table_schema(Table, Columns)
    ->  (   validate_values(Columns, Values)
        ->  get_next_id(Table, Id),
            assert(table_data(Table, Id, Values)),
            save_table_data(Table),
            Response = _{status: "success", message: "Record inserted", id: Id}
        ;   Response = _{status: "error", message: "Invalid values for table schema"}
        )
    ;   Response = _{status: "error", message: "Table does not exist"}
    ).

select_handler(Dict, Response) :-
    Table = Dict.get(table),
    Where = Dict.get(where, _{}),
    (   table_schema(Table, Columns)
    ->  findall(_{id: Id, data: Data}, 
                (table_data(Table, Id, Data), match_where(Data, Columns, Where)),
                Results),
        Response = _{status: "success", table: Table, columns: Columns, rows: Results}
    ;   Response = _{status: "error", message: "Table does not exist"}
    ).

update_handler(Dict, Response) :-
    Table = Dict.get(table),
    Set = Dict.get(set),
    Where = Dict.get(where, _{}),
    (   table_schema(Table, Columns)
    ->  findall(Id, 
                (table_data(Table, Id, Data), match_where(Data, Columns, Where)),
                Ids),
        length(Ids, Count),
        update_records(Table, Ids, Set, Columns),
        save_table_data(Table),
        Response = _{status: "success", message: "Records updated", count: Count}
    ;   Response = _{status: "error", message: "Table does not exist"}
    ).

delete_handler(Dict, Response) :-
    Table = Dict.get(table),
    Where = Dict.get(where, _{}),
    (   table_schema(Table, Columns)
    ->  findall(Id, 
                (table_data(Table, Id, Data), match_where(Data, Columns, Where)),
                Ids),
        length(Ids, Count),
        delete_records(Table, Ids),
        save_table_data(Table),
        Response = _{status: "success", message: "Records deleted", count: Count}
    ;   Response = _{status: "error", message: "Table does not exist"}
    ).

validate_values(Columns, Values) :-
    dict_keys(Values, ValueKeys),
    sort(ValueKeys, SortedKeys),
    sort(Columns, SortedCols),
    subset(SortedKeys, SortedCols).

get_next_id(Table, Id) :-
    findall(ExistingId, table_data(Table, ExistingId, _), Ids),
    (   Ids = []
    ->  Id = 1
    ;   max_list(Ids, MaxId),
        Id is MaxId + 1
    ).

match_where(_, _, Where) :-
    dict_keys(Where, []), !.
match_where(Data, Columns, Where) :-
    dict_pairs(Where, _, Pairs),
    forall(member(Key-Value, Pairs),
           (   nth0(Idx, Columns, Key),
               nth0(Idx, Data, Value)
           )).

update_records(_, [], _, _).
update_records(Table, [Id|Ids], Set, Columns) :-
    retract(table_data(Table, Id, OldData)),
    update_data(OldData, Set, Columns, NewData),
    assert(table_data(Table, Id, NewData)),
    update_records(Table, Ids, Set, Columns).

update_data(OldData, Set, Columns, NewData) :-
    dict_pairs(Set, _, Pairs),
    foldl(update_field(Columns), Pairs, OldData, NewData).

update_field(Columns, Key-Value, Data, UpdatedData) :-
    nth0(Idx, Columns, Key),
    replace_nth(Idx, Data, Value, UpdatedData).

replace_nth(0, [_|T], X, [X|T]) :- !.
replace_nth(N, [H|T], X, [H|R]) :-
    N > 0,
    N1 is N - 1,
    replace_nth(N1, T, X, R).

delete_records(_, []).
delete_records(Table, [Id|Ids]) :-
    retractall(table_data(Table, Id, _)),
    delete_records(Table, Ids).

save_schema(Table) :-
    db_directory(Dir),
    atom_concat(Dir, Table, BasePath),
    atom_concat(BasePath, '_schema.pl', FilePath),
    table_schema(Table, Columns),
    open(FilePath, write, Stream),
    format(Stream, ':- dynamic table_schema/2.~n', []),
    format(Stream, 'table_schema(~q, ~q).~n', [Table, Columns]),
    close(Stream).

save_table_data(Table) :-
    db_directory(Dir),
    atom_concat(Dir, Table, BasePath),
    atom_concat(BasePath, '_data.pl', FilePath),
    open(FilePath, write, Stream),
    format(Stream, ':- dynamic table_data/3.~n', []),
    forall(table_data(Table, Id, Data),
           format(Stream, 'table_data(~q, ~w, ~q).~n', [Table, Id, Data])),
    close(Stream).

load_all_tables :-
    db_directory(Dir),
    atom_concat(Dir, '*.pl', Pattern),
    expand_file_name(Pattern, Files),
    forall(member(File, Files), consult(File)).

:- initialization(init_db).

% Example usage (run these in Prolog console):
% ?- start_server.
% 
% Then use curl or any HTTP client:
% curl -X POST http://localhost:8080/query \
%   -H "Content-Type: application/json" \
%   -d '{"type":"create_table","table":"users","columns":["name","email","age"]}'
%
% curl -X POST http://localhost:8080/query \
%   -H "Content-Type: application/json" \
%   -d '{"type":"insert","table":"users","values":["John Doe","john@example.com",30]}'
%
% curl -X POST http://localhost:8080/query \
%   -H "Content-Type: application/json" \
%   -d '{"type":"select","table":"users","where":{"name":"John Doe"}}'
%
% curl -X POST http://localhost:8080/query \
%   -H "Content-Type: application/json" \
%   -d '{"type":"update","table":"users","set":{"age":31},"where":{"name":"John Doe"}}'
%
% curl -X POST http://localhost:8080/query \
%   -H "Content-Type: application/json" \
%   -d '{"type":"delete","table":"users","where":{"name":"John Doe"}}'
