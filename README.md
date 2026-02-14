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
  - `Credit`: int (Teaching Hours per Week / ชั่วโมงการสอนต่อสัปดาห์)
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
  - **Purpose**: Authenticate user and retrieve JWT token.
  - body: `{"username": "admin", "password": "123456"}`
- `get {{url}}/auth?search=&page=1&perpage=100` (Admin)
  - **Purpose**: List all users (teachers, students, admins).
  - body: `{}`
- `put {{url}}/auth/:id` (Admin)
  - **Purpose**: Update user credentials or role.
  - body: `{"username": "...", "human_type": "...", "role": ...}`
- `delete {{url}}/auth/:id?humantype=s` (Admin)
  - **Purpose**: Delete a user account.
  - body: `{}`

### Teacher (`/teacher`)
- `get {{url}}/teacher/:id` (Admin/User)
  - **Purpose**: Get detailed information about a specific teacher.
  - body: `{}`
- `post {{url}}/teacher` (Admin)
  - **Purpose**: Register a new teacher.
  - body: `{"name": "...", "resume": "...", "username": "...", "password": "...", "role": 1}`
- `get {{url}}/teacher?search=&page=1&perpage=100` (Admin/User)
  - **Purpose**: List all teachers with pagination and search.
  - body: `{}`
- `put {{url}}/teacher/:id` (Admin)
  - **Purpose**: Update teacher profile (name, resume).
  - body: `{"name": "...", "resume": "..."}`
- `delete {{url}}/teacher/:id` (Admin)
  - **Purpose**: Remove a teacher from the system.
  - body: `{}`
- `get {{url}}/teacher/:id/mysubject?search=&page=1&perpage=100` (Admin/User)
  - **Purpose**: List subjects a teacher is capable of teaching, with preference scores.
  - body: `{}`
- `post {{url}}/teacher/:id/mysubject` (Admin)
  - **Purpose**: Assign subjects to a teacher with a preference level (suitability).
  - body: `[{"subject_id": 1, "preference": 5}]`
- `delete {{url}}/teacher/:id/mysubject` (Admin)
  - **Purpose**: Unassign specific subjects from a teacher.
  - body: `{"ids": [1, 2]}`
- `delete {{url}}/teacher/:id/mysubject/all` (Admin)
  - **Purpose**: Remove all subject assignments for a specific teacher.
  - body: `{}`
- `post {{url}}/teacher/aievaluate` (Admin)
  - **Purpose**: AI-powered evaluation to suggest or match teachers to subjects.
  - body: `{}`
- `ws {{url}}/teacher/aievaluate/ws` (Admin)
  - **Purpose**: Real-time AI evaluation with progress updates via WebSocket.
  - **Behavior**: Streams progress text messages, ends with `"FINISHED"` or `"Error: ..."`.

### Student (`/student`)
- `get {{url}}/student?search=&page=1&perpage=100` (Admin/User)
  - **Purpose**: List all students.
  - body: `{}`
- `post {{url}}/student` (Admin)
  - **Purpose**: Register a new student and assign them to a classroom/curriculum.
  - body: `{"name": "...", "curriculum_id": 1, "year": 2024, "classroom_id": 1, "username": "...", "password": "...", "role": 1}`
- `get {{url}}/student/:id` (Admin/User)
  - **Purpose**: Get detailed profile of a student.
  - body: `{}`
- `put {{url}}/student/:id` (Admin)
  - **Purpose**: Update student details.
  - body: `{"name": "...", "curriculum_id": 1, "year": 2024, "classroom_id": 1}`
- `delete {{url}}/student/:id` (Admin)
  - **Purpose**: Remove a student from the system.
  - body: `{}`

### Subject (`/subject`)
- `get {{url}}/subject?search=&page=1&perpage=100` (Admin/User)
  - **Purpose**: List all available subjects in the system.
  - body: `{}`
- `post {{url}}/subject` (Admin)
  - **Purpose**: Create one or more new subjects.
  - body: `{"names": ["Math", "Science"]}`
- `get {{url}}/subject/:id` (Admin/User)
  - **Purpose**: Get details of a specific subject.
  - body: `{}`
- `put {{url}}/subject/:id` (Admin)
  - **Purpose**: Rename a subject.
  - body: `{"name": "Math"}`
- `delete {{url}}/subject/:id` (Admin)
  - **Purpose**: Delete a subject.
  - body: `{}`

### Classroom (`/classroom`)
- `get {{url}}/classroom?search=&page=1&perpage=100` (Admin/User)
  - **Purpose**: List all classrooms.
  - body: `{}`
- `post {{url}}/classroom` (Admin)
  - **Purpose**: Create a new classroom bound to a curriculum.
  - body: `{"name": "Room 101", "pre_curriculum_id": 1, "curriculum_id": 2, "year": 2024}`
- `get {{url}}/classroom/:id` (Admin/User)
  - **Purpose**: Get details of a specific classroom.
  - body: `{}`
- `put {{url}}/classroom/:id` (Admin)
  - **Purpose**: Update classroom details.
  - body: `{"name": "Room 101 Updated", "pre_curriculum_id": 1, "curriculum_id": 2, "year": 2024}`
- `delete {{url}}/classroom/:id` (Admin)
  - **Purpose**: Delete a classroom.
  - body: `{}`

### Term (`/term`)
- `get {{url}}/term` (Admin/User)
  - **Purpose**: List all academic terms (e.g., Term 1, Term 2).
  - body: `{}`
- `post {{url}}/term` (Admin)
  - **Purpose**: Create a new term.
  - body: `{"name": "Term 1"}`
- `put {{url}}/term/:id` (Admin)
  - **Purpose**: Rename a term.
  - body: `{"name": "Term 1 Updated"}`
- `delete {{url}}/term/:id` (Admin)
  - **Purpose**: Delete a term.
  - body: `{}`

### Curriculum (`/curriculum`)
- `get {{url}}/curriculum?search=&page=1&perpage=100` (Admin/User)
  - **Purpose**: List all curriculums.
  - body: `{}`
- `post {{url}}/curriculum` (Admin)
  - **Purpose**: Create a new curriculum.
  - body: `{"name": "...", "pre_curriculum_id": 1}`
- `get {{url}}/curriculum/:id` (Admin/User)
  - **Purpose**: Get curriculum details, including its subjects.
  - body: `{}`
- `put {{url}}/curriculum/:id` (Admin)
  - **Purpose**: Update curriculum name.
  - body: `{"name": "..."}`
- `delete {{url}}/curriculum/:id` (Admin)
  - **Purpose**: Delete a curriculum.
  - body: `{}`
- `post {{url}}/curriculum/:id/subject` (Admin)
  - **Purpose**: Manually add subjects to a curriculum.
  - body: `[1, 2, 3]`
- `delete {{url}}/curriculum/subject` (Admin)
  - **Purpose**: Remove subjects from a curriculum.
  - body: `{"ids": [1]}`
- `patch {{url}}/curiculum/subject/term` (Admin)
  - **Purpose**: Assign specific terms to subjects within a curriculum.
  - body: `[{"id": 1, "term_id": 1}, {"id": 2, "term_id": null}]`
- `post {{url}}/curriculum/map-precurriculum` (Admin)
  - **Purpose**: Automatically sync and populate subjects from PreCurriculum to all Curriculums.
  - body: `{}`

### PreCurriculum (`/precurriculum`)
- `get {{url}}/precurriculum?search=&page=1&perpage=100` (Admin/User)
  - **Purpose**: List all pre-curriculums (templates for curriculums).
  - body: `{}`
- `post {{url}}/precurriculum` (Admin)
  - **Purpose**: Create a new pre-curriculum (e.g., "Standard Science Course 2024").
  - body: `{"name": "..."}`
- `get {{url}}/precurriculum/:id` (Admin/User)
  - **Purpose**: Get details of a pre-curriculum.
  - body: `{}`
- `put {{url}}/precurriculum/:id` (Admin)
  - **Purpose**: Update pre-curriculum name.
  - body: `{"name": "..."}`
- `delete {{url}}/precurriculum/:id` (Admin)
  - **Purpose**: Delete a pre-curriculum.
  - body: `{}`
- `post {{url}}/precurriculum/:id/subject` (Admin)
  - **Purpose**: Define subjects and credits within a pre-curriculum.
  - body: `[{"subject_name": "Math", "credit": 3}]`
- `delete {{url}}/precurriculum/subject/:id` (Admin)
  - **Purpose**: Remove a subject from a pre-curriculum.
  - body: `{}`
- `post {{url}}/precurriculum/import` (Admin)
  - **Purpose**: Import a curriculum from a file (Text, CSV) or Image (Snapshot/Screenshot) using AI.
  - body (multipart): `file=@curriculum.png` (Multimodal support) OR `text="..."`
  - **Note**: AI will extract subject names and teaching hours (credit) automatically.
- `ws {{url}}/precurriculum/import/ws` (Admin)
  - **Purpose**: Real-time AI import with progress updates via WebSocket.
  - **Behavior**: 
    1. Connect.
    2. Client sends 1st message (Text or Base64 Image).
    3. Server streams progress messages.
    4. Ends with `"SUCCESS: [ID]"` or `"ERROR: ..."`.
- `post {{url}}/precurriculum/import/webhook` (Admin)
  - **Purpose**: Asynchronous AI import with Webhook callback.
  - **Body**: `{"callback_url": "https://...", "text": "...", "image_base64": "..."}`
  - **Response**: `202 Accepted` immediately.
  - **Behavior**: Server sends a POST to `callback_url` when finished with `{"status": "success", "id": ...}` or error details.

### Scadul Student (`/scadulstudent`)
- `get {{url}}/scadulstudent?search=&page=1&perpage=100` (Admin/User)
  - **Purpose**: List all generated student schedules.
  - body: `{}`
- `post {{url}}/scadulstudent` (Admin)
  - **Purpose**: Manually create a student schedule entry.
  - body: `{"classroom_id": 1, "use_in": "2024/1"}`
- `get {{url}}/scadulstudent/:id` (Admin/User)
  - **Purpose**: Get specific student schedule details.
  - body: `{}`
- `put {{url}}/scadulstudent/:id` (Admin)
  - **Purpose**: Update student schedule metadata.
  - body: `{"classroom_id": 1, "use_in": "2024/1"}`
- `delete {{url}}/scadulstudent/:id` (Admin)
  - **Purpose**: Delete a student schedule.
  - body: `{}`
- `post {{url}}/scadulstudent/:id/subject` (Admin)
  - **Purpose**: Add a specific subject slot (with teacher and order) to a student's schedule.
  - body: `[{"teacher_id": 1, "subject_id": 1, "order": 1}]`
- `delete {{url}}/scadulstudent/subject` (Admin)
  - **Purpose**: Remove a subject slot from a student's schedule.
  - body: `{"ids": [1]}`

### Scadul Teacher (`/scadulteacher`)
- `get {{url}}/scadulteacher?search=&page=1&perpage=100` (Admin/User)
  - **Purpose**: List all generated teacher schedules.
  - body: `{}`
- `post {{url}}/scadulteacher` (Admin)
  - **Purpose**: Manually create a teacher schedule entry.
  - body: `{"teacher_id": 1, "use_in": "2024/1"}`
- `get {{url}}/scadulteacher/:id` (Admin/User)
  - **Purpose**: Get specific teacher schedule details.
  - body: `{}`
- `put {{url}}/scadulteacher/:id` (Admin)
  - **Purpose**: Update teacher schedule metadata.
  - body: `{"teacher_id": 1, "use_in": "2024/1"}`
- `delete {{url}}/scadulteacher/:id` (Admin)
  - **Purpose**: Delete a teacher schedule.
  - body: `{}`
- `post {{url}}/scadulteacher/:id/subject` (Admin)
  - **Purpose**: Add a specific subject slot (teaching assignment) for a teacher.
  - body: `[{"teacher_id": 1, "subject_id": 1, "order": 1}]`
- `delete {{url}}/scadulteacher/subject` (Admin)
  - **Purpose**: Remove a teaching assignment slot.
  - body: `{"ids": [1]}`

### Schedule (`/schedule`)
- `post {{url}}/schedule/generate` (Admin)
  - **Purpose**: **Core Feature**. Triggers the automatic scheduling engine. This loads all data, runs the scheduling algorithm, and saves the resulting Student and Teacher schedules to the database.
  - **AI-Driven Mode**: Send an optional `prompt` in the body to specify constraints using natural language.
  - **Body**: 
    ```json
    {
      "prompt": "ครูสมชายไม่ว่างช่วงบ่ายวันพุธ และครูลีน่าอยากสอนคาบเช้าวันศุกร์ 2 คาบติด"
    }
    ```
  - **Time Period Mapping**:
    - The system uses a 1-40 period scale (8 periods/day, Mon-Fri).
    - Mon: 1-8, Tue: 9-16, Wed: 17-24, Thu: 25-32, Fri: 33-40.
  - **AI Logic**: 
    - The AI translates natural language into "Preferences".
    - `Preference: 10` (Must teach), `Preference: -10` (Cannot teach), etc.
    - These constraints are weighted against the core algorithm to ensure a valid, conflict-free schedule.
