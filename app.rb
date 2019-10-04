#! /usr/bin/env ruby

require 'bundler'
Bundler.require

class App < Sinatra::Base
  use Rack::Cors do
    allow do
      origins '*'
      resource '*', methods: %i[get post]
    end
  end

  get '/' do
    status 200
  end

  post '/' do
    # TODO: validation
    object = s3_bucket.object("#{SecureRandom.uuid}.png")
    # TODO: async
    object.upload_file(params[:file][:tempfile], acl: 'public-read')
    image_url = object.public_url

    status 201
    { image_url: image_url }.to_json
  end

  private

  def s3
    s3_params = {
      region: 'ap-northeast-1',
      credentials: Aws::Credentials.new(ENV['AWS_ACCESS_KEY_ID'], ENV['AWS_SECRET_ACCESS_KEY'])
    }
    @s3 ||= Aws::S3::Resource.new(s3_params)
  end

  def s3_bucket
    @s3_bucket ||= begin
                     s3_bucket = s3.bucket(ENV['S3_BUCKET'])
                     s3_bucket = s3.create_bucket(bucket: ENV['S3_BUCKET']) unless s3_bucket.exists?
                     s3_bucket
                   end
  end
end
