-- +goose Up
-- +goose StatementBegin
CREATE TABLE student (
     student_id BIGSERIAL PRIMARY KEY,
     student_name TEXT NOT NULL DEFAULT '',
     grade INT NOT NULL DEFAULT 0,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE class_info (
    id BIGSERIAL PRIMARY KEY,
    student_id INT NOT NULL DEFAULT 0,
    class_name TEXT DEFAULT '',
    CONSTRAINT fk_student FOREIGN KEY (student_id) REFERENCES student(student_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table class_info;
drop table student;
-- +goose StatementEnd
