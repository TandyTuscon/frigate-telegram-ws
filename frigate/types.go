package frigate

type EventStruct struct {
    EventID      string                 // Event unique identifier
    Camera       string                 // Camera name
    Label        string                 // Event label (e.g., person, dog)
    Timestamp    int64                  // Event timestamp
    Data         map[string]interface{} // Additional data about the event
    HasSnapshot  bool                   // Indicates if the event has a snapshot
    HasClip      bool                   // Indicates if the event has a video clip
    ID           string                 // Identifier for snapshots or clips
}
