package jobs

import "github.com/riverqueue/river"

const nonRetryableJobMaxAttempts = 1

func jobInsertOpts(queue string) *river.InsertOpts {
	return &river.InsertOpts{
		Queue:       queue,
		MaxAttempts: nonRetryableJobMaxAttempts,
	}
}

func jobInsertOptsWithUnique(queue string, unique river.UniqueOpts) *river.InsertOpts {
	opts := jobInsertOpts(queue)
	opts.UniqueOpts = unique
	return opts
}

func jobInsertOptsWithMetadataAndUnique(queue string, metadata []byte, unique river.UniqueOpts) *river.InsertOpts {
	opts := jobInsertOptsWithUnique(queue, unique)
	opts.Metadata = metadata
	return opts
}
