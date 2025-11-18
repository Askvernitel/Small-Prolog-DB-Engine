:- use_module(library(http/thread_httpd)).
:- use_module(library(http/http_dispatch)).
:- use_module(library(http/http_json)).
:- use_module(library(http/http_cors)).
:- use_module(library(http/json)).
:- dynamic table_data/3.

% HTTP route handlers
:- http_handler(root(query), handle_query, []).

% Start the server
start_server :-
    start_server(8080).

start_server(Port) :-
    http_server(http_dispatch, [port(Port)]),
    format('Server started on port ~w~n', [Port]),
    format('Data directory: ./data~n', []),
    ensure_data_directory.

% Ensure data directory exists
ensure_data_directory :-
    (   exists_directory('./data')
    ->  true
    ;   make_directory('./data')
    ).

% Main query handler - with error handling
handle_query(Request) :-
    cors_enable(Request, [methods([get,post,options])]),
    catch(
        (
            http_read_json_dict(Request, Dict),
            format(user_error, 'Received dict: ~w~n', [Dict]),
            process_query(Dict, Response),
            format(user_error, 'Response: ~w~n', [Response])
        ),
        Error,
        (
            format(user_error, 'Error processing query: ~w~n', [Error]),
            format(atom(ErrorMsg), '~w', [Error]),
            Response = _{status: "error", message: ErrorMsg}
        )
    ),
    reply_json_dict(Response).

% Process different query types
process_query(Dict, Response) :-
    (   get_dict(type, Dict, Type),
        get_dict(table, Dict, Table)
    ->  atom_string(TableAtom, Table),
        process_query_type(Type, TableAtom, Dict, Response)
    ;   Response = _{status: "error", message: "missing type or table field"}
    ).

% CREATE TABLE
process_query_type("create_table", Table, Dict, Response) :-
    (   get_dict(columns, Dict, Columns)
    ->  create_table(Table, Columns, Response)
    ;   Response = _{status: "error", message: "missing columns field"}
    ).

% INSERT
process_query_type("insert", Table, Dict, Response) :-
    (   get_dict(values, Dict, Values)
    ->  insert_row(Table, Values, Response)
    ;   Response = _{status: "error", message: "missing values field"}
    ).

% SELECT
process_query_type("select", Table, Dict, Response) :-
    (   get_dict(where, Dict, Where)
    ->  select_rows(Table, Where, Response)
    ;   select_rows(Table, _{}, Response)
    ).

% UPDATE
process_query_type("update", Table, Dict, Response) :-
    (   get_dict(set, Dict, Set)
    ->  (   get_dict(where, Dict, Where)
        ->  update_rows(Table, Set, Where, Response)
        ;   update_rows(Table, Set, _{}, Response)
        )
    ;   Response = _{status: "error", message: "missing set field"}
    ).

% DELETE
process_query_type("delete", Table, Dict, Response) :-
    (   get_dict(where, Dict, Where)
    ->  delete_rows(Table, Where, Response)
    ;   delete_rows(Table, _{}, Response)
    ).

% Unknown query type
process_query_type(Type, _Table, _Dict, Response) :-
    format(atom(Msg), 'unknown query type: ~w', [Type]),
    Response = _{status: "error", message: Msg}.

% CREATE TABLE implementation
create_table(Table, Columns, Response) :-
    table_file(Table, File),
    (   exists_file(File)
    ->  Response = _{status: "error", message: "table already exists"}
    ;   TableData = _{
            name: Table,
            columns: Columns,
            next_id: 1,
            rows: []
        },
        save_table(File, TableData),
        Response = _{
            status: "success",
            message: "table created",
            table: Table,
            columns: Columns
        }
    ).

% INSERT implementation
insert_row(Table, Values, Response) :-
    table_file(Table, File),
    (   load_table(File, TableData)
    ->  get_dict(columns, TableData, Columns),
        get_dict(next_id, TableData, ID),
        get_dict(rows, TableData, Rows),
        length(Columns, ColCount),
        length(Values, ValCount),
        (   ColCount =:= ValCount
        ->  create_row_data(Columns, Values, RowData),
            NewRow = _{id: ID, data: RowData},
            append(Rows, [NewRow], NewRows),
            NextID is ID + 1,
            put_dict(_{rows: NewRows, next_id: NextID}, TableData, UpdatedData),
            save_table(File, UpdatedData),
            Response = _{
                status: "success",
                message: "row inserted",
                table: Table,
                id: ID
            }
        ;   format(atom(Msg), 'expected ~w values, got ~w', [ColCount, ValCount]),
            Response = _{status: "error", message: Msg}
        )
    ;   Response = _{status: "error", message: "table does not exist"}
    ).

% SELECT implementation
select_rows(Table, Where, Response) :-
    table_file(Table, File),
    (   load_table(File, TableData)
    ->  get_dict(columns, TableData, Columns),
        get_dict(rows, TableData, Rows),
        filter_rows(Rows, Where, MatchedRows),
        length(MatchedRows, Count),
        Response = _{
            status: "success",
            table: Table,
            columns: Columns,
            rows: MatchedRows,
            count: Count
        }
    ;   Response = _{status: "error", message: "table does not exist"}
    ).

% UPDATE implementation
update_rows(Table, Set, Where, Response) :-
    table_file(Table, File),
    (   load_table(File, TableData)
    ->  get_dict(rows, TableData, Rows),
        update_matching_rows(Rows, Set, Where, UpdatedRows, Count),
        put_dict(rows, TableData, UpdatedRows, UpdatedData),
        save_table(File, UpdatedData),
        format(atom(Msg), '~w rows updated', [Count]),
        Response = _{
            status: "success",
            message: Msg,
            table: Table,
            count: Count
        }
    ;   Response = _{status: "error", message: "table does not exist"}
    ).

% DELETE implementation
delete_rows(Table, Where, Response) :-
    table_file(Table, File),
    (   load_table(File, TableData)
    ->  get_dict(rows, TableData, Rows),
        exclude(matches_where(Where), Rows, RemainingRows),
        length(Rows, OrigCount),
        length(RemainingRows, RemCount),
        Count is OrigCount - RemCount,
        put_dict(rows, TableData, RemainingRows, UpdatedData),
        save_table(File, UpdatedData),
        format(atom(Msg), '~w rows deleted', [Count]),
        Response = _{
            status: "success",
            message: Msg,
            table: Table,
            count: Count
        }
    ;   Response = _{status: "error", message: "table does not exist"}
    ).

% Helper predicates
table_file(Table, File) :-
    format(atom(File), './data/~w.json', [Table]).

save_table(File, Data) :-
    open(File, write, Stream, [encoding(utf8)]),
    json_write_dict(Stream, Data),
    nl(Stream),
    close(Stream).

load_table(File, Data) :-
    exists_file(File),
    open(File, read, Stream, [encoding(utf8)]),
    json_read_dict(Stream, Data),
    close(Stream).

create_row_data([], [], _{}).
create_row_data([Col|Cols], [Val|Vals], RowData) :-
    create_row_data(Cols, Vals, RestData),
    atom_string(ColAtom, Col),
    put_dict(ColAtom, RestData, Val, RowData).

filter_rows([], _, []).
filter_rows([Row|Rows], Where, [Row|Filtered]) :-
    matches_where(Where, Row),
    !,
    filter_rows(Rows, Where, Filtered).
filter_rows([_|Rows], Where, Filtered) :-
    filter_rows(Rows, Where, Filtered).

matches_where(Where, Row) :-
    get_dict(data, Row, Data),
    dict_pairs(Where, _, Pairs),
    all_match(Pairs, Data).

all_match([], _).
all_match([Key-Value|Rest], Data) :-
    get_dict(Key, Data, DataValue),
    DataValue = Value,
    all_match(Rest, Data).

update_matching_rows([], _, _, [], 0).
update_matching_rows([Row|Rows], Set, Where, [UpdatedRow|Updated], Count) :-
    matches_where(Where, Row),
    !,
    get_dict(data, Row, Data),
    merge_dicts(Set, Data, NewData),
    put_dict(data, Row, NewData, UpdatedRow),
    update_matching_rows(Rows, Set, Where, Updated, RestCount),
    Count is RestCount + 1.
update_matching_rows([Row|Rows], Set, Where, [Row|Updated], Count) :-
    update_matching_rows(Rows, Set, Where, Updated, Count).

merge_dicts(Source, Target, Merged) :-
    dict_pairs(Source, _, Pairs),
    apply_updates(Pairs, Target, Merged).

apply_updates([], Dict, Dict).
apply_updates([Key-Value|Rest], Dict, Result) :-
    put_dict(Key, Dict, Value, Updated),
    apply_updates(Rest, Updated, Result).

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