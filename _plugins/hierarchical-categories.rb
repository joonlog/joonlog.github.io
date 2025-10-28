# frozen_string_literal: true

module Jekyll
  # Generates hierarchical category pages
  # Primary categories: /categories/primary/
  # Secondary categories: /categories/primary/secondary/
  class HierarchicalCategoryPage < Page
    def initialize(site, base, primary, secondary = nil)
      @site = site
      @base = base
      @dir = if secondary
               File.join('categories', primary.downcase.gsub(' ', '-'))
             else
               'categories'
             end
      @name = if secondary
                "#{secondary.downcase.gsub(' ', '-')}.html"
              else
                "#{primary.downcase.gsub(' ', '-')}.html"
              end

      self.process(@name)
      self.read_yaml(File.join(base, '_layouts'), 'category.html')

      # Set category name for display
      self.data['title'] = secondary || primary
      self.data['category_name'] = secondary || primary
      self.data['primary_category'] = primary
      self.data['secondary_category'] = secondary

      # Filter posts based on category hierarchy
      if secondary
        # Secondary category: posts where categories[0] == primary AND categories[1] == secondary
        self.data['posts'] = site.posts.docs.select do |post|
          post.data['categories'] &&
          post.data['categories'][0] == primary &&
          post.data['categories'][1] == secondary
        end
      else
        # Primary category: posts where categories[0] == primary
        self.data['posts'] = site.posts.docs.select do |post|
          post.data['categories'] &&
          post.data['categories'][0] == primary
        end
      end
    end
  end

  class HierarchicalCategoryPageGenerator < Generator
    safe true
    priority :low

    def generate(site)
      return unless site.layouts.key? 'category'

      # Collect all primary and secondary categories
      primary_categories = {}

      site.posts.docs.each do |post|
        next unless post.data['categories']

        categories = post.data['categories']
        primary = categories[0]
        secondary = categories[1] if categories.length > 1

        next unless primary

        primary_categories[primary] ||= Set.new
        primary_categories[primary] << secondary if secondary
      end

      # Generate pages for each primary category
      primary_categories.each do |primary, secondaries|
        # Generate primary category page
        site.pages << HierarchicalCategoryPage.new(site, site.source, primary)

        # Generate secondary category pages
        secondaries.each do |secondary|
          next if secondary.nil?
          site.pages << HierarchicalCategoryPage.new(site, site.source, primary, secondary)
        end
      end
    end
  end
end
