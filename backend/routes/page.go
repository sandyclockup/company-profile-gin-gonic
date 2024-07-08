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

func PageRoutes() []RouteSource {
	routes := []RouteSource{
		{
			"/page/ping",
			"GET",
			false,
			controllers.PagePing,
		},
		{
			"/page/home",
			"GET",
			false,
			controllers.PageHome,
		},
		{
			"/page/about",
			"GET",
			false,
			controllers.PageAbout,
		},
		{
			"/page/service",
			"GET",
			false,
			controllers.PageService,
		},
		{
			"/page/faq",
			"GET",
			false,
			controllers.PageFaq,
		},
		{
			"/page/contact",
			"GET",
			false,
			controllers.PageContact,
		},
		{
			"/page/message",
			"POST",
			false,
			controllers.PageMessage,
		},
		{
			"/page/subscribe",
			"POST",
			false,
			controllers.PageSubscribe,
		},
		{
			"/page/uploads",
			"GET",
			false,
			controllers.PageGetFile,
		},
	}
	return routes
}
