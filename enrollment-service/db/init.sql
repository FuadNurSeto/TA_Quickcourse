CREATE TABLE enrollments (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(30) NOT NULL,
    course_id INT NOT NULL, -- Tanpa foreign key karena beda database
    status VARCHAR(15) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Penegakan Aturan Bisnis BR-02 (Satu pendaftaran aktif per mahasiswa per course)
CREATE UNIQUE INDEX satu_pendaftaran_aktif 
ON enrollments (student_id, course_id) 
WHERE status = 'ACTIVE';