/**
 * This file is part of the Sandy Andryanto Company Profile Website.
 *
 * @author     Sandy Andryanto <sandy.andryanto.dev@gmail.com>
 * @copyright  2024
 *
 * For the full copyright and license information,
 * please view the LICENSE.md file that was distributed
 * with this source code.
 */
 
 package routes

 import "backend/controllers"
 
 func ArticleRoutes() []RouteSource {
	 routes := []RouteSource{
		 {
			 "/article/list",
			 "GET",
			 false,
			 controllers.ArticleList,
		 },
		 {
			"/article/detail/:slug",
			"GET",
			false,
			controllers.ArticleDetail,
		},
		{
			"/article/comments/:id",
			"GET",
			false,
			controllers.ArticleCommentList,
		},
		{
			"/article/comments/:id",
			"POST",
			true,
			controllers.ArticleCommentCreate,
		},
	 }
	 return routes
 }