$(function(){
    $("#form").submit(function(){
        var keyword = $(this).find("[name=keyword]").val();
        $.getJSON("/search/" + keyword, function(data){
            console.log(data)
            //search-result
            
        })
    })
});