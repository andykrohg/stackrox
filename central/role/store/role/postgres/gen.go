package postgres

//go:generate pg-table-bindings-wrapper --type=storage.Role --migration-seq 46 --migrate-from rocksdb
