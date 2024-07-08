package api

type CreateWishlistParams struct {
	Name         string `json:"name" binding:"required"`
	OptionUserID string `json:"option_user_id" binding:"required"`
}

type WishlistItem struct {
	Name           string `json:"name"`
	WishlistID     string `json:"wishlist_id"`
	WishlistItemID string `json:"wishlist_item_id"`
	OptionUserID   string `json:"option_user_id"`
	MainImage     string `json:"main_image"`
}

type ListWishlistRes struct {
	List    []WishlistItem `json:"list"`
	IsEmpty bool           `json:"is_empty"`
}

type RemoveWishlistParams struct {
	WishlistID string `json:"wishlist_id" binding:"required"`
}

type RemoveWishlistRes struct {
	Success    bool   `json:"success"`
	WishlistID string `json:"wishlist_id" binding:"required"`
}

type RemoveWishlistItemParams struct {
	WishlistItemID string `json:"wishlist_item_id" binding:"required"`
}

type RemoveWishlistItemRes struct {
	Success      bool   `json:"success"`
	OptionUserID string `json:"option_user_id" binding:"required"`
}

type CreateWishlistItemParams struct {
	WishlistID   string `json:"wishlist_id" binding:"required"`
	OptionUserID string `json:"option_user_id" binding:"required"`
}

type WishlistOffsetParams struct {
	OptionOffset   int    `json:"option_offset"`
	MainOptionType string `json:"main_option_type" binding:"required"`
	WishlistID     string `json:"wishlist_id" binding:"required"`
	Currency       string `json:"currency" binding:"required"`
}


type ListExperienceWishlistOptionRes struct {
	List         []ExperienceOptionData `json:"list"`
	OptionOffset int                    `json:"option_offset"`
	OnLastIndex  bool                   `json:"on_last_index"`
	WishlistID   string                 `json:"wishlist_id"`
}

type ListExperienceWishlistEventRes struct {
	List         []ExperienceEventData `json:"list"`
	OptionOffset int                    `json:"option_offset"`
	OnLastIndex  bool                   `json:"on_last_index"`
	WishlistID   string                 `json:"wishlist_id"`
}


