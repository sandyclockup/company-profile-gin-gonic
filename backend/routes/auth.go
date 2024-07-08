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
 
 func AuthRoutes() []RouteSource {
	 routes := []RouteSource{
		 {
			 "/auth/login",
			 "POST",
			 false,
			 controllers.AuthLogin,
		 },
		 {
			"/auth/register",
			"POST",
			false,
			controllers.AuthRegister,
		},
		{
			"/auth/confirm/:token",
			"GET",
			false,
			controllers.AuthConfirm,
		},
		{
			"/auth/email/forgot",
			"POST",
			false,
			controllers.AuthEmailForgot,
		},
		{
			"/auth/email/reset/:token",
			"POST",
			false,
			controllers.AuthEmailReset,
		},
	 }
	 return routes
 }