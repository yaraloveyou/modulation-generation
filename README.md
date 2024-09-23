# Проект на Go для генерации Java Spring CRUD
```(Проект на стадии разработки)```
## Описание проекта

Данный проект представляет собой генератор CRUD приложения на Java Spring на основе структуры таблиц, описанных в JSON файле. Программа на Go считывает JSON файл конфигураций, где указаны таблицы, их поля и связи между ними, а также настройки файловой системы для создания соответствующих сущностей, репозиториев, контроллеров и сервисов в Java.

## Особенности

- Чтение конфигурации проекта из JSON файла.
- Автоматическая генерация Java Spring классов для CRUD операций.
- Поддержка создания сущностей с Foreign Key связями.
- Настройка файловой системы для размещения сгенерированных файлов.

### Пример структуры JSON файла:

```json
{
    "Name": "library_system",
    "Lombok": "true",
    "Tables": [
        {
            "name": "User",
            "modifier": "protected",
            "fields": {
                "id": "int;identity",
                "name": "string",
                "email": "string",
                "roleId": "int;foreign_key{Role};unique",
                "borrowedBookId": "int;foreign_key{Book}"
            }
        },
        {
            "name": "Role",
            "modifier": "protected",
            "fields": {
                "id": "int;identity",
                "name": "string"
            }
        }
    ]
}
```
## Как это работает
- Чтение настроек из JSON: Программа на Go загружает JSON файл, содержащий описание таблиц и их полей.
- Генерация Java сущностей: На основе описания таблиц, программа создает Java классы сущностей для каждой таблицы, с учетом аннотаций для работы с базой данных.
- Создание репозиториев, сервисов и контроллеров: По конфигурации, указанной в JSON, программа настраивает пакеты и генерирует необходимые файлы для работы с CRUD операциями в проекте на Java Spring.

## Пример сгенерированных Java сущностей
Сущность ```User```

```java
package generated_sources.test.test_project.domain;

@Getter
@Setter
@ToString
@NoArgsConstructor
@AllArgsConstructor
protected class User implements Serializable {
	@Id
	@GeneratedValue(strategy = GenerationType.IDENTITY)
	@Column(name = "id")
	private Integer id;

	@Column(name = "name")
	private String name;

	@Column(name = "email")
	private String email;

	@OneToMany(cascade = CascadeType.ALL, mappedBy = "user")
	private List<Role> role;

}
```
Сущность ```Role```

```java
package generated_sources.test.test_project.domain;

@Getter
@Setter
@ToString
@NoArgsConstructor
@AllArgsConstructor
protected class Role implements Serializable {
	@Id
	@GeneratedValue(strategy = GenerationType.IDENTITY)
	@Column(name = "id")
	private Integer id;

	@Column(name = "name")
	private String name;

	@ManyToOne(cascade = CascadeType.ALL)
	private User user;

}
```
## Настройки файловой системы
JSON файл также содержит информацию о том, как должны быть организованы папки для размещения сгенерированных файлов.

Пример структуры папок:
```
{
    "package": "test.test_project",
    "folders": {
        "entity": "domain",
        "repository": "repository",
        "controller": "controller",
        "service": "service.dto"
    }
}
```
