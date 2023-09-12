namespace go chat

struct CreateChatReq {
    1: required i64 user_id     (api.vd = "$ > 0", api.query = "user_id")
    2: required i64 session_id  (api.query = "session_id")
}

struct CreateChatResp {
}

struct CreateReplyReq {
    1: required i64 user_id     (api.vd = "$ > 0")
    2: required i64 session_id  (api.vd = "$ > 0")
    3: required string msg      (api.vd = "$ != ''")
}

struct CreateReplyResp {
    1: required Message message
}

struct GetSessionListReq {
    1: required i64 user_id
}

struct GetSessionListResp {
    1: required list<i64> sessions
}

struct Message {
    1: i64 id
    2: i64 session_id
    3: i64 user_id
    4: string data
    5: i64 create_time
    6: i64 update_time
    7: string audio_key
}

struct GetSessionMsgReq {
    1: required i64 user_id
    2: required i64 session_id
}

struct GetSessionMsgResp {
    1: required list<Message> messages
}


service ChatService {
    CreateChatResp CreateChat(1: CreateChatReq request) (api.get="/chat");
    GetSessionListResp GetSessionList(1: GetSessionListReq request) (api.get="/session")
    GetSessionMsgResp GetSessionMsg(1: GetSessionMsgReq request) (api.get="/message")
    CreateReplyResp CreateReply(1: CreateReplyReq request) (api.post="/reply")
}

