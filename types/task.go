// Types to be shared over the module.
package types

// TaskSettings settings for the task.
type TaskSettings struct {
	Users               []int64 `bson:"users"`
	Url                 string  `bson:"url"`
	CurrentTextTemplate string  `bson:"current_text"`
}
