# 使用详情

## 代码生成命令

`go run -mod=mod github.com/99designs/gqlgen generate --verbose`

## 测试方式
1. 运行服务
2. 打开localhost:8080
3. 创建和查询

``` GraphQL
mutation createTodo {
  createTodo(input: {text: "todo", userId: "2"}) {
    user {
      id
    }
    text
    done
  }
}

query findTodos {
  todos {
    text
    done
    user {
      name
    }
  }
}

```