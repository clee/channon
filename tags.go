package main

import (
	"log"
	"net/http"

	"github.com/mholt/binding"
	"github.com/unrolled/render"
	"github.com/zenazn/goji/web"
)

func (t *TagCrumb) FieldMap(r *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&t.TagName: binding.Field{
			Form:     "tagName",
			Required: true,
		},
	}
}

func tagIndex(tags []*Tag, tag *Tag) int {
	for index, t := range tags {
		if *t == *tag {
			return index
		}
	}
	return -1
}

func listTagsHandler(pm *PlanManager) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		ren := render.New(render.Options{})
		ren.JSON(w, http.StatusOK, pm.tags)
	}
}

func addTagHandler(pm *PlanManager) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		t := new(TagCrumb)
		errs := binding.Bind(r, t)
		if errs.Handle(w) {
			return
		}

		log.Printf("adding tag named: %s\n", t.TagName)
		tag := Tag(t.TagName)

		go func() {
			pm.AddTag(&tag)
			pm.lock <- 0
		}()
		<-pm.lock

		ren := render.New(render.Options{})
		ren.JSON(w, http.StatusOK, map[string]string{"tagAdded": t.TagName})
	}
}

func deleteTagHandler(pm *PlanManager) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		t := new(TagCrumb)
		errs := binding.Bind(r, t)
		if errs.Handle(w) {
			return
		}

		tag := Tag(t.TagName)

		go func() {
			pm.DeleteTag(&tag)
			pm.lock <- 0
		}()
		<-pm.lock

		ren := render.New(render.Options{})
		ren.JSON(w, http.StatusOK, map[string]string{"tagDeleted": t.TagName})
	}
}
