-- database/schema.sql
CREATE TABLE IF NOT EXISTS public."VideoJobs"
(
    id text COLLATE pg_catalog."default" NOT NULL,
    "uploadId" text COLLATE pg_catalog."default" NOT NULL,
    index integer,
    CONSTRAINT "VideoJobs_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."HlsJobs"
(
    id text COLLATE pg_catalog."default" NOT NULL,
    "uploadId" text COLLATE pg_catalog."default" NOT NULL,
    status boolean NOT NULL DEFAULT false,
    CONSTRAINT "HlsJobs_pkey" PRIMARY KEY (id)
);