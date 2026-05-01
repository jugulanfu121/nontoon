-- database/query.sql

-- name: GetAllVideoJobs :many
SELECT id, "uploadId", index
	FROM public."VideoJobs";

-- name: AddVideoJob :exec
INSERT INTO public."VideoJobs"(
	id, "uploadId", index)
	VALUES ($1, $2, $3);

-- name: GetLatestUploadedChunk :one
SELECT id, "uploadId", index 
FROM public."VideoJobs" 
WHERE "uploadId" = $1
ORDER BY index DESC
LIMIT 1;

-- name: AddHlsJob :exec
INSERT INTO public."HlsJobs"(
	id, "uploadId", status)
	VALUES ($1, $2, $3);

-- name: UpdateHlsJobStatus :exec
UPDATE public."HlsJobs"
	SET status=$1
	WHERE "uploadId"=$2;