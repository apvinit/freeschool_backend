# Freeschool API

Api backend for the Freeschool Application

## Api Endpoints

### Categories

``` json
// Category JSON Object
{
  "id" : 1,
  "title" : "Math",
  "iconURL" : "https://location/to/some/icon"
}

```

**POST** `/api/categories` 

**GET** `/api/categories` 

**GET** `/api/categories/:id` 

**PUT** `/api/categories/:id` 

**DELETE** `/api/categories/:id` 

### Courses

``` json
// Course JSON object
{
  "id" : 1,
  "title": "Go Programming",
  "description": "Programming With Golang",
  "categoryID" : 12,
  "iconURL" : "https://location/to/some/icon"
}
```

**POST** `/api/courses` 

**GET** `/api/courses` 

**GET** `/api/courses/:id` 

**PUT** `/api/courses/:id` 

**DELETE** `/api/courses/:id` 

### Modules

``` json
// Topic JSON Object
{
  "id" : 1,
  "title": "Go Basics",
  "description" : "Basics of go programming",
}
```

**POST** `/api/modules` 

**GET** `/api/modules` 

**GET** `/api/modules/:id` 

**PUT** `/api/modules/:id` 

**DELETE** `/api/modules/:id` 

### Lessons

``` json
// Lesson JSON Object
{
  "id" : 1,
  "title": "Go Basics",
  "courseID" : 1,
  "contents" : [] // array of [Content] Objects
}
```

**POST** `/api/lessons` 

**GET** `/api/lessons` 

**GET** `/api/lessons/:id` 

**PUT** `/api/lessons/:id` 

**DELETE** `/api/lessons/:id` 

### Contents

``` json
{
  "id" : 1,
  "title" : "Intro to variables",
  "description":"Variables data types",
  "lessonID" : 12,
  "contentType" : "Video | Text | PDF",
  "data" :  "https://url/to/video"
}
```

**POST** `/api/contents` 

**GET** `/api/contents` 

**GET** `/api/contents/:id` 

**PUT** `/api/contents/:id` 

**DELETE** `/api/contents/:id` 
