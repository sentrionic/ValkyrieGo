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
	BanYourselfError       = "you cannot ban yourself"
	KickYourselfError      = "you cannot kick yourself"
	UnbanYourselfError     = "you cannot unban yourself"
	OneChannelRequired     = "A server needs at least one channel"
	ChannelLimitError      = "channel limit is 50"
	DMYourselfError        = "you cannot dm yourself"
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
	Unauthorized   = "Not Authorized"
)

// Message Errors
const (
	MessageOrFileRequired    = "Either a message or a file is required"
	MessageEmptyError        = "Message must not be empty"
	InvalidMimeType          = "file must be of type 'image' or 'audio'"
	EditMessageError         = "Only the author can edit the message"
	InvalidRequestParameters = "Invalid request parameters. See errors"
	TextRequiredError        = "Text is required"
	DeleteMessageError       = "Only the author or owner can delete the message"
	DeleteDMMessageError     = "Only the author can delete the message"
)
