# doorkeeper

    Http Сервис исполняющий запросы к другим ресурсам.
    
    Config:
        -Address - По умолчанию(0.0.0.0:8000)
        -Workers - колличество рутин для worker. По умолчания 2.
        -TaskCountOnPage - количество просьб на одной странице.
    
    Хендлеры:
       - GET /v1/fetchtask?task={"method":"GET","address":"http://example.com"}
            - task:
                - method - разрешенно передавать get, post, put, delete.
                - adderss - сслыка на внешний ресурс.
            Response:
                json {id, status, headers, response_length}
       - GET /v1/tasks/{page}
            - page - номер страницы
            Response:
                json tasks
       - DELETE /v1/tasks/{id}
            - id - уникальный номер просьбы.
            Response:
                status.