# A simple example on usage of go-on-rails generator

This is an simple example for the [go-on-rails](https://github.com/goonr/go-on-rails), a Rails generator.

You can take the example as a tutorial, too. I'll make it as simple and clear as possible to show how to use the go-on-rails generator to generate Golang codes in a Rails app.

### Environments

* macOS Sierra v10.12.4
* Ruby v2.3.3
* Rails v5.0.2
* MySQL v5.7.11
* Golang v1.7.4 darwin/amd64

### Build a Rails app

Firstly, we will follow the tutorial in Rails guides, to [build a (very) simple weblog](http://guides.rubyonrails.org/getting_started.html). We may not copy that whole steps, but the models mainly.

Let's create a new Rails app:

```bash
rails new example_simple --api --database mysql --skip-bundle
```

change to the new directory, add the gem go-on-rails:

```bash
# edit Gemfile
gem 'go-on-rails', '~> 0.0.9'
```
and then bundle:

```bash
bundle install
```

### Create some models

We'll build two models: article and comment, here's an article has_many comments association.

```bash
rails g model Article title:string text:text
```

```bash
rails g model Comment commenter:string body:text article_id:integer
```

You'd better add some restrict to the migration files, eg. add `null: false` restriction to the `title` column:

```ruby
# the migration file under db/migrate
class CreateArticles < ActiveRecord::Migration[5.0]
  def change
    create_table :articles do |t|
      t.string :title, null: false
      t.text :text

      t.timestamps
    end
  end
end
```

And meanwhile we add some presence and length validations to the models:

```ruby
# app/models/article.rb
class Article < ApplicationRecord
  has_many :comments, dependent: :destroy

  validates :title, presence: true, length: { in: 10..30 }
  validates :text, presence: true, length: { minimum: 20 }
end

# app/models/comment.rb
class Comment < ApplicationRecord
  belongs_to :article

  validates :commenter, presence: true
  validates :body, presence: true, length: { minimum: 20 }
end
```

### Generate Golang codes

Now it's time to show magic!

Run the command:

```bash
rails g gor dev
```
and a new directory named `go_app` with some files will be made under the root of Rails app:

```
# directory tree of go_app

├── controllers
│   └── home_controller.go
├── main.go
├── models
│   ├── db.go
│   ├── gor_article.go
│   └── gor_comment.go
├── public
│   └── favicon.ico
└── views
    └── index.tmpl
```

then let's install Golang dependancies(some packages) we'll use:

```bash
rails gor:deps
```

You can see completed alert after some while, and make sure you now under the `go_app`

```bash
go run main.go
```
