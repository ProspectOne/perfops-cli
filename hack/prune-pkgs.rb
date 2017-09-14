#!/usr/bin/env ruby

require 'json'
require 'pp'
require 'rest-client'
require 'time'

LIMIT = 2

API_TOKEN = ENV['PACKAGECLOUD_TOKEN']
USER = 'p1'
REPOSITORY = 'perfops'

def get_pkgs(base_url, pkgs_url)
	url = base_url + pkgs_url
	pkgs = RestClient.get(url)
	return JSON.parse(pkgs)
end

def get_versions(base_url, version_url)
	url = base_url + version_url
	begin
		versions = RestClient.get(url)
	rescue => e
		puts "Skipping #{version_url}: #{e.message}"
		return []
	end

	parsed = JSON.parse(versions)
	return parsed.sort_by { |x| Time.parse(x["created_at"]) }
end

def delete_version(base_url, version)
	url = base_url + version['destroy_url']
	puts "Deleting #{version['destroy_url']}"
	begin
		RestClient.delete(url)
	rescue => e
		puts "ERROR: #{e.message}"
	end
end

def prune_pkgs(base_url, pkgs_url)
	pkgs = get_pkgs(base_url, pkgs_url)
	i = 0
	while i < pkgs.length do
		pkg = pkgs[i]
		if pkg['versions_count'] > LIMIT then
			versions = get_versions(base_url, pkg['versions_url'])
			j = 0
			max = versions.length - LIMIT
			while j < max do
				delete_version(base_url, versions[j])
				j += 1
			end
		end
		i += 1
	end
end

base_url = "https://#{API_TOKEN}:@packagecloud.io"

["deb", "rpm"].each { |t|
	pkgs_url = "/api/v1/repos/#{USER}/#{REPOSITORY}/packages/#{t}.json"
	prune_pkgs(base_url, pkgs_url)
}
