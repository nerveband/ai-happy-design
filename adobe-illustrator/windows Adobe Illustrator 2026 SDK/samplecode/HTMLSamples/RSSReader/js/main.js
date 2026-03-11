var http = require('http');
var FeedParser = require('feedparser');
var Mustache = require('mustache');

var feedUrl = "http://feeds.bbci.co.uk/news/rss.xml";
//var feedUrl = "http://blogs.adobe.com/indesignsdk/feed/";

$(function() {		
	// load feed..
	fetchFeed();
});

function fetchFeed() {
	$.mobile.loading('show');

	var items = [];
	http.get(feedUrl, function(res) {

		if (res.statusCode == '200')
		{			
			res.pipe(new FeedParser({}))
	            .on('error', function(error){
	                alert("FeedParser Error");
	            })
	            .on('meta', function(meta) {
					// not needed
	            })
	            .on('readable', function() {
	                var stream = this, item;
	                while (item = stream.read()) {
	                    items.push(item);
	                }  
	            })
	            .on('end', function() {				
					$("#item_list").empty();
					var $homePage = $("#home").clone();
											
					$.each(items, function(i, item) {
						// add item to list of items..
						$("#item_list").append($('<li></li>').html('<a href="#detail_' + i + '">' + item.title + '</a>'));
					
		            	// create a detail page for the item
		            	$.get('../templates/detail_page.mustache', function(template) { 
						
						    var templateData = {
								id: "detail_" + i,
						    	title: item.title,
						    	date: item.pubDate,
						    	description: item.description,
								link: item.link
						    };

						    $("body").append(Mustache.render($(template).html(), templateData));
						});
						
						$("#item_list").listview("refresh");
					});   
	            });

        } else {
        	alert("Bad HTTP status code: " + res.statusCode);
        }

		$.mobile.loading('hide');

    }).on('error', function(e) {
		$.mobile.loading('hide');
		alert("HTTP Error: " + e.message);
	});	
}