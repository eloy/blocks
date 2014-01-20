package blocks_test

import (
	"github.com/harlock/blocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


var _ = Describe("Route", func() {
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
		rootRoute, homeRoute, barRoute *blocks.Route
	)


	BeforeEach(func() {
		blocks.R.Reset()
		rootRoute = blocks.R.Root(HomeController{}, "Index")
		homeRoute = blocks.R.Get("/home", FooController{}, "Home")
		barRoute = blocks.R.Namespace("api").Get("bar/:name", BarController{}, "Bar")
	})


	Describe("ControllerName()", func() {
		It("Should return the the controller name lowercase", func() {
			Expect(rootRoute.ControllerName()).To(Equal("home"))
			Expect(homeRoute.ControllerName()).To(Equal("foo"))
			Expect(barRoute.ControllerName()).To(Equal("bar"))
		})
	})

	Describe("ActionName()", func() {
		It("Should return the the controller name lowercase", func() {
			Expect(rootRoute.ActionName()).To(Equal("index"))
			Expect(homeRoute.ActionName()).To(Equal("home"))
			Expect(barRoute.ActionName()).To(Equal("bar"))
		})
	})

})
