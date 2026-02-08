## Migration

```bash
go run infra/gormDB/migration/migration.go   
```

## Run

```bash
go run main.go
```

## API Endpoints Summary

### Auth (`/auth`)
- `POST {{url}}/auth/login` - Login
  - Body: `{"username": "admin", "password": "password"}`
- `GET {{url}}/auth` - List Users (Admin)
- `PUT {{url}}/auth/:id` - Update User (Admin)
  - Body: `{"username": "...", "human_type": "...", "role": ...}`
- `DELETE {{url}}/auth/:id` - Delete User (Admin)
  - Query: `?humantype=s` (s: student, t: teacher)

### Teacher (`/teacher`)
- `GET {{url}}/teacher/:id` - Get Teacher (Admin/User)
- `POST {{url}}/teacher` - Create Teacher (Admin)
  - Body: `{"name": "...", "resume": "...", "username": "...", "password": "...", "role": 1}`
- `GET {{url}}/teacher` - List Teachers (Admin/User)
- `PUT {{url}}/teacher/:id` - Update Teacher (Admin)
  - Body: `{"name": "...", "resume": "..."}`
- `DELETE {{url}}/teacher/:id` - Delete Teacher (Admin)
- `GET {{url}}/teacher/:id/mysubject` - View My Subjects (Admin/User)
- `POST {{url}}/teacher/:id/mysubject` - Add Subjects (Admin)
  - Body: `[{"subject_id": 1, "preference": 5}]`
- `DELETE {{url}}/teacher/:id/mysubject` - Remove Subjects (Admin)
  - Body: `{"ids": [1, 2]}`
- `POST {{url}}/teacher/aieverlute` - AI Evaluation (Admin)

### Student (`/student`)
- `GET {{url}}/student` - List Students (Admin/User)
- `POST {{url}}/student` - Create Student (Admin)
  - Body: `{"name": "...", "curriculum_id": 1, "year": 2024, "classroom_id": 1, "username": "...", "password": "...", "role": 1}`
- `GET {{url}}/student/:id` - Get Student (Admin/User)
- `PUT {{url}}/student/:id` - Update Student (Admin)
  - Body: `{"name": "...", "curriculum_id": 1, "year": 2024, "classroom_id": 1}`
- `DELETE {{url}}/student/:id` - Delete Student (Admin)

### Subject (`/subject`)
- `GET {{url}}/subject` - List Subjects (Admin/User)
- `POST {{url}}/subject` - Create Subjects (Admin)
  - Body: `{"names": ["Math", "Science"]}`
- `GET {{url}}/subject/:id` - Get Subject (Admin/User)
- `PUT {{url}}/subject/:id` - Update Subject (Admin)
  - Body: `{"name": "Math"}`
- `DELETE {{url}}/subject/:id` - Delete Subject (Admin)

### Classroom (`/classroom`)
- `GET {{url}}/classroom` - List Classrooms (Admin/User)
  - Search: `ByName`, `ByPreCurriculumName`, `ByCurriculumName`, `ByYear`
- `POST {{url}}/classroom` - Create Classroom (Admin)
  - Body: `{"name": "Room 101", "pre_curriculum_id": 1, "curriculum_id": 2, "year": 2024}`
- `GET {{url}}/classroom/:id` - Get Classroom (Admin/User)
- `PUT {{url}}/classroom/:id` - Update Classroom (Admin)
  - Body: `{"name": "Room 101 Updated", "pre_curriculum_id": 1, "curriculum_id": 2, "year": 2024}`
- `DELETE {{url}}/classroom/:id` - Delete Classroom (Admin)

### Term (`/term`)
- `GET {{url}}/term` - List Terms (Admin/User)
- `POST {{url}}/term` - Create Term (Admin)
  - Body: `{"name": "Term 1"}`
- `GET {{url}}/term/:id` - Get Term (Admin/User)
- `PUT {{url}}/term/:id` - Update Term (Admin)
  - Body: `{"name": "Term 1 Updated"}`
- `DELETE {{url}}/term/:id` - Delete Term (Admin)

### Curriculum (`/curriculum`)
- `GET {{url}}/curriculum` - List Curriculums (Admin/User)
- `POST {{url}}/curriculum` - Create Curriculum (Admin)
  - Body: `{"name": "..."}`
- `GET {{url}}/curriculum/:id` - Get Curriculum (Admin/User)
- `PUT {{url}}/curriculum/:id` - Update Curriculum (Admin)
  - Body: `{"name": "..."}`
- `DELETE {{url}}/curriculum/:id` - Delete Curriculum (Admin)
- `POST {{url}}/curriculum/:id/subject` - Add Subject to Curriculum (Admin)
  - Body: `[{"subject_in_pre_curriculum_id": 1, "term_id": 1}]`
- `DELETE {{url}}/curriculum/subject` - Remove Subjects (Admin)
  - Body: `{"ids": [1]}`

### PreCurriculum (`/precurriculum`)
- `GET {{url}}/precurriculum` - List PreCurriculums (Admin/User)
- `POST {{url}}/precurriculum` - Create PreCurriculum (Admin)
  - Body: `{"name": "..."}`
- `GET {{url}}/precurriculum/:id` - Get PreCurriculum (Admin/User)
- `PUT {{url}}/precurriculum/:id` - Update PreCurriculum (Admin)
  - Body: `{"name": "..."}`
- `DELETE {{url}}/precurriculum/:id` - Delete PreCurriculum (Admin)
- `POST {{url}}/precurriculum/:id/subject` - Add Subject (Admin)
  - Body: `[{"subject_name": "Math", "credit": 3}]`
- `DELETE {{url}}/precurriculum/subject/:id` - Remove Subject (Admin)

### Scadul Student (`/scadulstudent`)
- `GET {{url}}/scadulstudent` - List (Admin/User)
- `POST {{url}}/scadulstudent` - Create (Admin)
  - Body: `{"classroom_id": 1, "use_in": "2024/1"}`
- `GET {{url}}/scadulstudent/:id` - Get (Admin/User)
- `PUT {{url}}/scadulstudent/:id` - Update (Admin)
  - Body: `{"classroom_id": 1, "use_in": "2024/1"}`
- `DELETE {{url}}/scadulstudent/:id` - Delete (Admin)
- `POST {{url}}/scadulstudent/:id/subject` - Add Subjects (Admin)
  - Body: `[{"teacher_id": 1, "subject_id": 1, "order": 1}]`
- `DELETE {{url}}/scadulstudent/subject` - Remove Subjects (Admin)
  - Body: `{"ids": [1]}`

### Scadul Teacher (`/scadulteacher`)
- `GET {{url}}/scadulteacher` - List (Admin/User)
- `POST {{url}}/scadulteacher` - Create (Admin)
  - Body: `{"teacher_id": 1, "use_in": "2024/1"}`
- `GET {{url}}/scadulteacher/:id` - Get (Admin/User)
- `PUT {{url}}/scadulteacher/:id` - Update (Admin)
  - Body: `{"teacher_id": 1, "use_in": "2024/1"}`
- `DELETE {{url}}/scadulteacher/:id` - Delete (Admin)
- `POST {{url}}/scadulteacher/:id/subject` - Add Subjects (Admin)
  - Body: `[{"teacher_id": 1, "subject_id": 1, "order": 1}]`
- `DELETE {{url}}/scadulteacher/subject` - Remove Subjects (Admin)
  - Body: `{"ids": [1]}`
