package enum;

// Status for a key value pair
type Status int32

const (
    // Modify means Status is updating
    Modify      Status = 0
    // Shared means the data is sharing
    Shared      Status = 1
    // Invalid means data is invalid
    Invalid     Status = 2
)