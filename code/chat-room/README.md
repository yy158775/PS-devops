项目名称：chat-room聊天室

server：服务端

client：客户端

原理：


chat-service:

提供微服务接口：


rpc Register(RegisterRequest) returns (RegisterResponse);

rpc Login(LoginRequest)       returns (LoginResponse);

rpc SendMessage(ChatMessage)  returns (Empty);