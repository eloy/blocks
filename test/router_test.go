package blocks_test

import (
	"github.com/harlock/blocks"
	m "github.com/harlock/blocks/test/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func newPather(method string, path string) blocks.Pather {
	p := m.NewMockPather(m.NewCtrl())
	p.EXPECT().Path().Return(path).AnyTimes()
	p.EXPECT().Method().Return(method).AnyTimes()
	return p
}


var _ = Describe("Router", func() {
	type HomeController struct {
		blocks.ApplicationController
	}

	type FooController struct {
		blocks.ApplicationController
	}

	type BarController struct {
		blocks.ApplicationController
	}

	type WadusController struct {
		blocks.ApplicationController
	}

	var (
		rootRoute, homeRoute, barRoute blocks.Routable
	)

	BeforeEach(func() {
		blocks.R.Reset()
		rootRoute = blocks.R.Root(HomeController{}, "Index")
		homeRoute = blocks.R.Get("/home", HomeController{}, "Home")
		barRoute = blocks.R.Namespace("api").Get("bar/:name", BarController{}, "Bar")
	})



	// Path()
	//--------------------------------------------------------------------

	Describe("Getting the Path method", func() {
		Context("Without parents", func() {
			It("Should return the controller name", func() {
				Expect(homeRoute.Path()).To(Equal("/home"))
			})
		})

		Context("With parents", func() {
			It("Should return the controller name", func() {
				Expect(barRoute.Path()).To(Equal("/api/bar/:name"))
			})
		})
	})



	// Match(Pathable)
	//--------------------------------------------------------------------
	Describe("Match()", func() {
		It("Should return true if the given path match the route", func() {

			// Route /
			path := newPather("GET", "/")
			Expect(homeRoute.Match(path)).ToNot(BeTrue())
			Expect(barRoute.Match(path)).ToNot(BeTrue())

			// Route /home
			path = newPather("GET", "/home")
			Expect(homeRoute.Match(path)).To(BeTrue())
			Expect(barRoute.Match(path)).ToNot(BeTrue())

			// Route /api/bar/test
			path = newPather("GET", "/api/bar/test")
			Expect(homeRoute.Match(path)).ToNot(BeTrue())
			Expect(barRoute.Match(path)).To(BeTrue())
		})
	})

	// Find(Pathable)
	//--------------------------------------------------------------------
	Describe("Find()", func() {

		BeforeEach(func() {
			wadus := blocks.R.Resources(WadusController{})
			wadus.Member().Resources(FooController{})
		})

		It("Should return the route configured for the root path", func() {
			path := newPather("GET", "/")
			route, found := blocks.R.Find(path)
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("home"))
			Expect(route.ActionName()).To(Equal("index"))
		})


		It("Should work with simple rules", func() {
			path := newPather("GET", "/home")
			route, found := blocks.R.Find(path)
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("home"))
			Expect(route.ActionName()).To(Equal("home"))
		})

		It("Should work with namespaces", func() {
			path := newPather("GET", "/api/bar/test")
			route, found := blocks.R.Find(path)
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("bar"))
			Expect(route.ActionName()).To(Equal("bar"))
		})

		It("should manage complex namespaces", func() {
			blocks.R.Root(HomeController{}, "Index")
			api := blocks.R.Namespace("api")
			api.Resources(HomeController{})
			apiAdmin := api.Namespace("admin")
			apiAdmin.Resources(BarController{})

			path := newPather("GET", "/api/admin/bar/")

			route, found := blocks.R.Find(path)

			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("bar"))
			Expect(route.ActionName()).To(Equal("index"))
		})

		It("Should work with resources", func() {
			blocks.R.Resources(WadusController{})
			route, found := blocks.R.Find(newPather("GET", "/wadus"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("wadus"))
			Expect(route.ActionName()).To(Equal("index"))

			route, found = blocks.R.Find(newPather("GET", "/wadus/3"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("wadus"))
			Expect(route.ActionName()).To(Equal("show"))

			route, found = blocks.R.Find(newPather("GET", "/wadus/3/edit"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("wadus"))
			Expect(route.ActionName()).To(Equal("edit"))

			route, found = blocks.R.Find(newPather("GET", "/wadus/new"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("wadus"))
			Expect(route.ActionName()).To(Equal("new"))

			route, found = blocks.R.Find(newPather("POST", "/wadus"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("wadus"))
			Expect(route.ActionName()).To(Equal("create"))

			route, found = blocks.R.Find(newPather("PUT", "/wadus/3"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("wadus"))
			Expect(route.ActionName()).To(Equal("update"))

			route, found = blocks.R.Find(newPather("DELETE", "/wadus/3"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("wadus"))
			Expect(route.ActionName()).To(Equal("destroy"))
		})

		It("Should work with resources of resources members", func() {
			blocks.R.Reset()
			wadus := blocks.R.Resources(WadusController{})
			wadus.Member().Resources(FooController{})


			route, found := blocks.R.Find(newPather("GET", "/wadus/1/foo"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("foo"))
			Expect(route.ActionName()).To(Equal("index"))

			route, found = blocks.R.Find(newPather("GET", "/wadus/1/foo/3"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("foo"))
			Expect(route.ActionName()).To(Equal("show"))

			route, found = blocks.R.Find(newPather("GET", "/wadus/1/foo/3/edit"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("foo"))
			Expect(route.ActionName()).To(Equal("edit"))

			route, found = blocks.R.Find(newPather("GET", "/wadus/1/foo/new"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("foo"))
			Expect(route.ActionName()).To(Equal("new"))

			route, found = blocks.R.Find(newPather("POST", "/wadus/1/foo"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("foo"))
			Expect(route.ActionName()).To(Equal("create"))

			route, found = blocks.R.Find(newPather("PUT", "/wadus/1/foo/3"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("foo"))
			Expect(route.ActionName()).To(Equal("update"))

			route, found = blocks.R.Find(newPather("DELETE", "/wadus/1/foo/3"))
			Expect(found).To(BeTrue())
			Expect(route.ControllerName()).To(Equal("foo"))
			Expect(route.ActionName()).To(Equal("destroy"))
		})


	})

})
