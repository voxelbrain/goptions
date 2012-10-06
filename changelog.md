# 2.1.0

# New features

* `goptions.Verbs` is of type `string` and will have selected verb name as value
  after parsing.

# 2.0.0


## Breaking changes

* Disallow multiple flag names for one member
* Remove `accumulate` option in favor of generic array support


## New features

* Add convenience function `ParseAndFail` to make common usage of the library
  a one-liner (see `readme_example.go`)
* Add a `Marshaler` interface to enable thrid-party types
* Add support for slices (and thereby for mutiple flag definitions)


## Minor changes

* Refactoring to get more flexibility
* Make a flag's default value accessible in the template context
