#include "server.h"

// 定义静态全局 app 对象，确保它的生命周期覆盖整个程序

uWS::App &setupRoutes()
{
        static uWS::App *app = nullptr; // 使用 static 保证生命周期

        if (!app)
        {
                app = new uWS::App(); // 动态分配内存来初始化 app
                // 设置路由
                app->get("/", [](auto *res, auto *req)
                         {
                                 res->end("Hello world!"); // 返回响应
                         });
        }
        else
        {
                printf("app is already initialized\n");
        }

        return *app;
}
