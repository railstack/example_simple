# A simple example on how to use go-on-rails generator

Of course, you can take this example as a tutorial, too. I'll make it as simple and clear as possible to show how to use the go-on-rails generator to generate Golang codes.

### Ruby and Rails versions

We choose the Ruby 2.3.3 and Rails 5.0.2 in the example, and all these're specified in the Gemfile, and we take MySQL as our database.

### Build a Rails app

Firstly, we will follow the tutorial in Rails guides, to [build a (very) simple weblog](http://guides.rubyonrails.org/getting_started.html). We may not copy that whole steps, but the models mainly.

Let's create a new Rails app:

```bash
rails new example_simple --api --database mysql --skip-bundle
```

change to the new directory, add the gem go-on-rails:

```bash
# edit Gemfile
gem 'go-on-rails', '>~0.0.9'
```
and then bundle:

```bash
bundle install
```

### Create some models

We'll build two models: Article and Comment, which have a has_many association.

