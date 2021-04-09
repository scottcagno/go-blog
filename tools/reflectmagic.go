package tools

import (
	"fmt"
	"reflect"
)

type Model struct {
	Name   string
	Fields []ModelField
	Count  int
}

type ModelField struct {
	Name  string
	Type  string
	Tag   string
	Value interface{}
}

func MakeModel(v interface{}, tag string) (*Model, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		// uncomment line below if you wish to only receive ptr values
		//return nil, fmt.Errorf("expected pointer, got: %v\n", val.Kind())
	}
	val = reflect.Indirect(val)
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got: %v\n", val.Kind())
	}
	typ := val.Type()
	if typ.NumField() < 1 {
		return nil, fmt.Errorf("expected non empty struct, got struct with %d fields\n", typ.NumField())
	}
	model := &Model{
		Name:   typ.String(),
		Fields: make([]ModelField, typ.NumField()),
		Count:  typ.NumField(),
	}
	for i := 0; i < typ.NumField(); i++ {
		fld := typ.Field(i)
		if !fld.Anonymous {
			model.Fields[i] = ModelField{
				Name:  fld.Name,
				Type:  fld.Type.Name(),
				Tag:   fld.Tag.Get(tag),
				Value: val.Field(i).Interface(),
			}
		}
	}
	return model, nil
}

var bsMain = `<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <!-- Bootstrap CSS -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-eOJMYsd53ii+scO/bJGFsiCZc+5NDVN2yr8+0RDqr0Ql0h+rP48ckxlpbzKgwra6" crossorigin="anonymous">

    <title>Hello, world!</title>
  </head>
  <body>
    <h1>Hello, world!</h1>

    <!-- Optional JavaScript; choose one of the two! -->

    <!-- Option 1: Bootstrap Bundle with Popper -->
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta3/dist/js/bootstrap.bundle.min.js" integrity="sha384-JEW9xMcG8R+pH31jmWH6WWP0WintQrMb4s7ZOdauHnUtxwoG2vI5DkLtS3qm9Ekf" crossorigin="anonymous"></script>

    <!-- Option 2: Separate Popper and Bootstrap JS -->
    <!--
    <script src="https://cdn.jsdelivr.net/npm/@popperjs/core@2.9.1/dist/umd/popper.min.js" integrity="sha384-SR1sx49pcuLnqZUnnPwx6FCym0wLsk5JZuNx2bPPENzswTNFaQU1RDvt3wT4gWFG" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta3/dist/js/bootstrap.min.js" integrity="sha384-j0CNLUeiqtyaRmlzUHCPZ+Gy5fQu0dQ6eZ/xAww941Ai1SxSY+0EQqNXNE6DZiVc" crossorigin="anonymous"></script>
    -->
  </body>
</html>`

var bsForm1 = `<!-- beg: sample form #1 -->
<div class="container">
  <br>
  <legend>Form #1, Normal Labels & Helpers</legend>
  <hr>
  <form action="#">
    
    <div class="mb-3">
      <label for="emailAddr1" class="form-label">Email address</label>
      <input type="email" class="form-control" id="emailAddr1" aria-describedby="emailAddr1-help">
      <div id="emailAddr1-help" class="form-text">Please login using your email address</div>
    </div>
    
    <div class="mb-3">
      <label for="pass1" class="form-label">Password</label>
      <input type="password" class="form-control" id="pass1" aria-describedby="pass1-help">
      <div id="pass1-help" class="form-text">Password must have at least 6 characters</div>
    </div>
    
    <div class="mb-3">
      <label for="textarea1" class="form-label">Comments</label>
      <textarea class="form-control" id="textarea1" aria-describedby="textarea1-help" style="height: 100px"></textarea>
      <div id="textarea1-help" class="form-text">Please leave a comment and let us know what you think</div>
    </div>

    <div class="mb-3">
    <label for="select1" class="form-label">Make a selection</label>
      <select id="select1" class="form-select" aria-label="Default select example" aria-describedby="select1-help">
        <option></option>
        <option value="1">One</option>
        <option value="2">Two</option>
        <option value="3">Three</option>
      </select>
      <div id="select1-help" class="form-text">No wrong selections here folks!</div>
    </div>

    <div class="mb-3">
      <input class="form-check-input" type="checkbox" id="remember1" aria-describedby="remember1-help">
      <label class="form-check-label" for="remember1">Remember me 1</label>
      <div id="remember1-help" class="form-text">Check this to have your login saved</div>
    </div>
    
    <div class="d-grid gap-2 d-md-flex justify-content-md-end">
      <button type="submit" class="btn btn-primary me-md-2">Save Form 1</button>
      <button class="btn btn-secondary">Cancel Form 1</button>
    </div>
    
  </form>
</div>
<!-- end: sample form #1 -->`

var bsForm2 = `<!-- beg: sample form #2 -->
<div class="container">
  <br>
  <legend>Form #2, Floating Labels & No Helpers</legend>
  <hr>
  <form action="#">

    <div class="form-floating mb-3">
      <input type="email" class="form-control" id="emailAddr2" placeholder="name@example.com">
      <label for="emailAddr2">Email address</label>
    </div>

    <div class="form-floating mb-3">
      <input type="password" class="form-control" id="pass2" placeholder="Password">
      <label for="pass2">Password</label>
    </div>

    <div class="form-floating mb-3">
      <textarea class="form-control" placeholder="Leave a comment here" id="textarea2" style="height: 100px"></textarea>
      <label for="textarea2">Comments</label>
    </div>

    <div class="mb-3">
      <select id="select2" class="form-select">
        <option selected>Make a selection</option>
        <option value="1">One</option>
        <option value="2">Two</option>
        <option value="3">Three</option>
      </select>
    </div>

    <div class="mb-3">
      <input class="form-check-input" type="checkbox" id="remember2">
      <label class="form-check-label" for="remember2">Remember me 2</label>
    </div>

    <div class="d-grid gap-2 d-md-flex justify-content-md-end">
      <button type="submit" class="btn btn-primary me-md-2">Save Form 2</button>
      <button class="btn btn-secondary">Cancel Form 2</button>
    </div>

  </form>
</div>
<!-- end: sample form #2 -->`

var bsForm3 = `<!-- beg: sample form #3 -->
<div class="container">
  <br>
  <legend>Form #3, Placeholders & No Helpers</legend>
  <hr>
  <form action="#">
    
    <div class="mb-3">
      <input type="email" class="form-control" id="emailAddr3" placeholder="Email address">
    </div>
    
    <div class="mb-3">
      <input type="password" class="form-control" id="pass3" placeholder="Password">
    </div>
    
    <div class="mb-3">
      <textarea class="form-control" id="textarea3" style="height: 100px" placeholder="Comments"></textarea>
    </div>
    
    <div class="mb-3">
      <select class="form-select" aria-label="Default select example">
        <option selected>Make a selection</option>
        <option value="1">One</option>
        <option value="2">Two</option>
        <option value="3">Three</option>
      </select>
    </div>
    
    <div class="mb-3">
      <input class="form-check-input" type="checkbox" id="remember3">
      <label class="form-check-label" for="remember3">Remember me 3</label>
    </div>
    
    <div class="d-grid gap-2 d-md-flex justify-content-md-end">
      <button type="submit" class="btn btn-primary me-md-2">Save Form 3</button>
      <button class="btn btn-secondary">Cancel Form 3</button>
    </div>
    
  </form>
</div>
<!-- end: sample form #3 -->`
