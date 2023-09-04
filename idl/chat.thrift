namespace go chat

struct CreateChatReq {
}

struct CreateChatResp {
}


service ChatService {
    CreateChatResp CreateChat(1: CreateChatReq request) (api.get="/chat");
}

