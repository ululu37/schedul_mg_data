## Migration
```bash
docker compose up -d
```

```bash
go run infra/gormDB/migration/migration.go   
```
```bash
insert into auths (username ,"password" ,"role" ,human_type) values ('admin','123456',0,'t');
```
## Run

```bash
go run main.go
```



## Domain Entities

### Core Entities
- **Auth**:
  - `ID`: uint (Primary Key)
  - `Username`: string (Unique)
  - `Password`: string
  - `HumanType`: string ('t': Teacher, 's': Student, '': Other)
  - `Role`: int (0: User, 1: Admin)
- **Teacher**:
  - `ID`: uint (Primary Key)
  - `Name`: string
  - `Resume`: text
  - `AuthID`: uint (Foreign Key to Auth)
- **Student**:
  - `ID`: uint (Primary Key)
  - `AuthID`: uint (Foreign Key to Auth)
  - `Name`: string
  - `CurriculumID`: *uint (Foreign Key to Curriculum)
  - `Year`: int
  - `ClassroomID`: *uint (Foreign Key to Classroom)
- **Subject**:
  - `ID`: uint (Primary Key)
  - `Name`: string
- **Term**:
  - `ID`: uint (Primary Key)
  - `Name`: string
- **Classroom**:
  - `ID`: uint (Primary Key)
  - `Name`: string
  - `PreCurriculumID`: *uint (Foreign Key to PreCurriculum)
  - `CurriculumID`: *uint (Foreign Key to Curriculum)
  - `Year`: *int

### Curriculum Management
- **PreCurriculum**:
  - `ID`: uint (Primary Key)
  - `Name`: string
- **SubjectInPreCurriculum**:
  - `ID`: uint (Primary Key)
  - `PreCurriculumID`: uint (Foreign Key to PreCurriculum)
  - `SubjectID`: uint (Foreign Key to Subject)
  - `Credit`: int
- **Curriculum**:
  - `ID`: uint (Primary Key)
  - `Name`: string
  - `PreCurriculumID`: uint (Foreign Key to PreCurriculum)
- **SubjectInCurriculum**:
  - `ID`: uint (Primary Key)
  - `CurriculumID`: uint (Foreign Key to Curriculum)
  - `SubjectInPreCurriculumID`: uint (Foreign Key to SubjectInPreCurriculum)
  - `TermID`: *uint (Foreign Key to Term)

### Scheduling
- **ScadulTeacher**:
  - `ID`: uint (Primary Key)
  - `TeacherID`: uint (Foreign Key to Teacher)
  - `UseIn`: string (YYYY/term)
- **SubjectInScadulTeacher**:
  - `ID`: uint (Primary Key)
  - `ScadulTeacherID`: uint (Foreign Key to ScadulTeacher)
  - `TeacherID`: uint (Foreign Key to Teacher)
  - `SubjectID`: uint (Foreign Key to Subject)
  - `Order`: int
- **ScadulStudent**:
  - `ID`: uint (Primary Key)
  - `ClassroomID`: uint (Foreign Key to Classroom)
  - `UseIn`: string (YYYY/term)
- **SubjectInScadulStudent**:
  - `ID`: uint (Primary Key)
  - `ScadulStudentID`: uint (Foreign Key to ScadulStudent)
  - `TeacherID`: uint (Foreign Key to Teacher)
  - `SubjectID`: uint (Foreign Key to Subject)
  - `Order`: int

### Teacher Management
- **TeacherMySubject**:
  - `ID`: uint (Primary Key)
  - `TeacherID`: uint (Foreign Key to Teacher)
  - `SubjectID`: uint (Foreign Key to Subject)
  - `Preference`: int

## API Endpoints Summary

### Auth (`/auth`)
- `post {{url}}/auth/login`
  - body: `{"username": "admin", "password": "password"}`
- `get {{url}}/auth?search=&page=1&perpage=100` (Admin)
  - body: `{}`
- `put {{url}}/auth/:id` (Admin)
  - body: `{"username": "...", "human_type": "...", "role": ...}`
- `delete {{url}}/auth/:id?humantype=s` (Admin)
  - body: `{}`

### Teacher (`/teacher`)
- `get {{url}}/teacher/:id` (Admin/User)
  - body: `{}`
- `post {{url}}/teacher` (Admin)
  - body: `{"name": "...", "resume": "...", "username": "...", "password": "...", "role": 1}`
- `get {{url}}/teacher?search=&page=1&perpage=100` (Admin/User)
  - body: `{}`
- `put {{url}}/teacher/:id` (Admin)
  - body: `{"name": "...", "resume": "..."}`
- `delete {{url}}/teacher/:id` (Admin)
  - body: `{}`
- `get {{url}}/teacher/:id/mysubject?search=&page=1&perpage=100` (Admin/User)
  - body: `{}`
- `post {{url}}/teacher/:id/mysubject` (Admin)
  - body: `[{"subject_id": 1, "preference": 5}]`
- `delete {{url}}/teacher/:id/mysubject` (Admin)
  - body: `{"ids": [1, 2]}`
- `post {{url}}/teacher/aieverlute` (Admin)
  - body: `{}`

### Student (`/student`)
- `get {{url}}/student?search=&page=1&perpage=100` (Admin/User)
  - body: `{}`
- `post {{url}}/student` (Admin)
  - body: `{"name": "...", "curriculum_id": 1, "year": 2024, "classroom_id": 1, "username": "...", "password": "...", "role": 1}`
- `get {{url}}/student/:id` (Admin/User)
  - body: `{}`
- `put {{url}}/student/:id` (Admin)
  - body: `{"name": "...", "curriculum_id": 1, "year": 2024, "classroom_id": 1}`
- `delete {{url}}/student/:id` (Admin)
  - body: `{}`

### Subject (`/subject`)
- `get {{url}}/subject?search=&page=1&perpage=100` (Admin/User)
  - body: `{}`
- `post {{url}}/subject` (Admin)
  - body: `{"names": ["Math", "Science"]}`
- `get {{url}}/subject/:id` (Admin/User)
  - body: `{}`
- `put {{url}}/subject/:id` (Admin)
  - body: `{"name": "Math"}`
- `delete {{url}}/subject/:id` (Admin)
  - body: `{}`

### Classroom (`/classroom`)
- `get {{url}}/classroom?search=&page=1&perpage=100` (Admin/User)
  - body: `{}`
- `post {{url}}/classroom` (Admin)
  - body: `{"name": "Room 101", "pre_curriculum_id": 1, "curriculum_id": 2, "year": 2024}`
- `get {{url}}/classroom/:id` (Admin/User)
  - body: `{}`
- `put {{url}}/classroom/:id` (Admin)
  - body: `{"name": "Room 101 Updated", "pre_curriculum_id": 1, "curriculum_id": 2, "year": 2024}`
- `delete {{url}}/classroom/:id` (Admin)
  - body: `{}`

### Term (`/term`)
- `get {{url}}/term` (Admin/User)
  - body: `{}`
- `post {{url}}/term` (Admin)
  - body: `{"name": "Term 1"}`
- `put {{url}}/term/:id` (Admin)
  - body: `{"name": "Term 1 Updated"}`
- `delete {{url}}/term/:id` (Admin)
  - body: `{}`

### Curriculum (`/curriculum`)
- `get {{url}}/curriculum?search=&page=1&perpage=100` (Admin/User)
  - body: `{}`
- `post {{url}}/curriculum` (Admin)
  - body: `{"name": "...", "pre_curriculum_id": 1}`
- `get {{url}}/curriculum/:id` (Admin/User)
  - body: `{}`
- `put {{url}}/curriculum/:id` (Admin)
  - body: `{"name": "..."}`
- `delete {{url}}/curriculum/:id` (Admin)
  - body: `{}`
- `post {{url}}/curriculum/:id/subject` (Admin)
  - body: `[1, 2, 3]`
- `delete {{url}}/curriculum/subject` (Admin)
  - body: `{"ids": [1]}`
- `patch {{url}}/curiculum/subject/term` (Admin)
  - body: `[{"id": 1, "term_id": 1}, {"id": 2, "term_id": null}]`

### PreCurriculum (`/precurriculum`)
- `get {{url}}/precurriculum?search=&page=1&perpage=100` (Admin/User)
  - body: `{}`
- `post {{url}}/precurriculum` (Admin)
  - body: `{"name": "..."}`
- `get {{url}}/precurriculum/:id` (Admin/User)
  - body: `{}`
- `put {{url}}/precurriculum/:id` (Admin)
  - body: `{"name": "..."}`
- `delete {{url}}/precurriculum/:id` (Admin)
  - body: `{}`
- `post {{url}}/precurriculum/:id/subject` (Admin)
  - body: `[{"subject_name": "Math", "credit": 3}]`
- `delete {{url}}/precurriculum/subject/:id` (Admin)
  - body: `{}`

### Scadul Student (`/scadulstudent`)
- `get {{url}}/scadulstudent?search=&page=1&perpage=100` (Admin/User)
  - body: `{}`
- `post {{url}}/scadulstudent` (Admin)
  - body: `{"classroom_id": 1, "use_in": "2024/1"}`
- `get {{url}}/scadulstudent/:id` (Admin/User)
  - body: `{}`
- `put {{url}}/scadulstudent/:id` (Admin)
  - body: `{"classroom_id": 1, "use_in": "2024/1"}`
- `delete {{url}}/scadulstudent/:id` (Admin)
  - body: `{}`
- `post {{url}}/scadulstudent/:id/subject` (Admin)
  - body: `[{"teacher_id": 1, "subject_id": 1, "order": 1}]`
- `delete {{url}}/scadulstudent/subject` (Admin)
  - body: `{"ids": [1]}`

### Scadul Teacher (`/scadulteacher`)
- `get {{url}}/scadulteacher?search=&page=1&perpage=100` (Admin/User)
  - body: `{}`
- `post {{url}}/scadulteacher` (Admin)
  - body: `{"teacher_id": 1, "use_in": "2024/1"}`
- `get {{url}}/scadulteacher/:id` (Admin/User)
  - body: `{}`
- `put {{url}}/scadulteacher/:id` (Admin)
  - body: `{"teacher_id": 1, "use_in": "2024/1"}`
- `delete {{url}}/scadulteacher/:id` (Admin)
  - body: `{}`
- `post {{url}}/scadulteacher/:id/subject` (Admin)
  - body: `[{"teacher_id": 1, "subject_id": 1, "order": 1}]`
- `delete {{url}}/scadulteacher/subject` (Admin)
  - body: `{"ids": [1]}`
