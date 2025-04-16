#include "server.h"
#include <iostream>

int main()
{
    // 设置路由并启动监听
    setupRoutes().listen(9021, [](auto *token)
                         {
        if (token) {
            std::cout << "Listening on port " << 9001 << std::endl;
        } })
        .run(); // 启动服务器并运行

    return 0;
}
