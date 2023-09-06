hertz:
	hz update -idl idl/chat.thrift

gorm:
	gentool -dsn "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local" -tables "message" -outPath ./dal/query

