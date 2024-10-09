package constants

// 处理结果
type HandlerResult int

const (
	NoHandlerResult     HandlerResult = iota + 1 //未处理
	PassHandlerResult                            // 通过
	RefuseHandlerResult                          // 拒绝
	CancelHandlerResult                          // 取消
)

// 群成员等级 1.创建者 2.管理者 3.普通成员
type GroupRoleLevel int

const (
	CreatorGroupRoleLevel GroupRoleLevel = iota + 1
	ManagerGroupRoleLevel
	AtLargeGroupRoleLevel
)

// 进群申请的方式： 1.邀请， 2.申请
type GroupJoinSource int

const (
	InviteGroupJoinSource GroupJoinSource = iota + 1
	PutInGroupJoinSource
)

const (
	ErrFriendReqAlreadyPassed = 1
)
