package user

type UserAddress struct {
	AddressID		string
	Alias				string
	Line1 			string
	Line2 			string
	Line3 			string
	City   			string
	State  			string
	ZipCode 		string
	Country 		string
}

type UserPreferences struct {
	Language string
	Timezone string
}

type UserMetadata struct {
	CreatedAt 		string
	UpdatedAt 		string
	LastLogin 		string
	EmailVerified bool
	MFAEnabled    bool
}

type ProfileSize string
const (
	ProfileSizeFull    ProfileSize = "full"    // Complete profile with address, preferences, metadata
	ProfileSizeMinimal ProfileSize = "minimal" // Only essential fields (user_id, email, name)
	ProfileSizeBasic   ProfileSize = "basic"   // Includes name fields but no address/preferences
)

type UserProfileGetRequest struct {
	UserID      string
	Size        ProfileSize
	BearerToken string
}

type UserProfileGetResponseBasic struct {
	UserID       	string
	FirstName 		string
	LastName  		string
}

type UserProfileGetResponseMinimal struct {
	UserID       	string
	Email    			string
	Phone					string
	FirstName 		string
	LastName  		string
	PasswordHash 	string
}

type UserProfileGetResponseFull struct {
	UserID        string
	Username      string
	Email         string
	Phone         string
	FirstName     string
	MiddleInitial string
	LastName      string
	PasswordHash  string
	Address       UserAddress
	Preferences   UserPreferences
	Metadata      UserMetadata
}

type UserProfileUpdateRequest struct {
	UserID        string
	Username      string
	Email         string
	Phone         string
	FirstName     string
	MiddleInitial string
	LastName      string
	Address       UserAddress
	Preferences   UserPreferences
	Metadata      UserMetadata
	BearerToken   string
}
