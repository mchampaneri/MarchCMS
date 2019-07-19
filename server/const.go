package main

// User role constants
const (
	adminUser    = 0 // Has rights to perform all opertaions
	designerUser = 1 // Can set theme, generate & arrange menu
	editorUser   = 2 // Can write page, write post, manage writer user, edit post written by writer user
	writerUser   = 3 // can write page, write post
)

// User acconunt status constant
const (
	inactiveAccount     = 0 // Newly created and unverified account - user yet to do his first login
	activeAccount       = 1 // verified and active account - user can login
	tempDisabledAccount = 2 // account temporary disabled
	permDisabledAccount = 3 // account permanantly disabled
)
