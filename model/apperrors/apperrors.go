package apperrors

// Guild Errors
const (
	NotAMember             = "not a member"
	AlreadyMember          = "already a member"
	GuildLimitReached      = "guild limit is 100"
	MustBeOwner            = "must be the owner for that"
	InvalidImageType       = "imageFile must be 'image/jpeg' or 'image/png'"
	MustBeMemberInvite     = "must be a member to fetch an invite"
	IsPermanentError       = "isPermanent is not a boolean"
	InvalidateInvitesError = "only the owner can invalidate invites"
	InvalidInviteError     = "Invalid Link or the server got deleted"
	BannedFromServer       = "You are banned from this server"
	DeleteGuildError       = "only the owner can delete their server"
	OwnerCantLeave         = "the owner cannot leave their server"
)

// Account Errors
const (
	DuplicateEmail      = "email already in use"
	PasswordsDoNotMatch = "passwords do not match"
)

// Friend Errors
const (
	AddYourselfError    = "You cannot add yourself"
	RemoveYourselfError = "You cannot remove yourself"
	AcceptYourselfError = "You cannot accept yourself"
	CancelYourselfError = "You cannot cancel yourself"
	UnableAddError      = "Unable to add user as friend"
	UnableRemoveError   = "Unable to remove the user"
	UnableAcceptError   = "Unable to accept the request"
)

// Generic Errors
const (
	InvalidSession = "provided session is invalid"
	ServerError    = "Something went wrong. Try again later"
)
