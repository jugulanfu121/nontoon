-- database/query.sql

-- name: GetAllVideoJobs :many

SELECT id, "uploadId", index
	FROM public."VideoJobs";

-- name: AddVideoJob :exec
INSERT INTO public."VideoJobs"(
	id, "uploadId", index)
	VALUES ($1, $2, $3);