$(function(){
    var searchResult = $('#search-result'), timeCost = $('#time-cost'), currentPage = 1;
    $("#form").submit(function(){
        if($(this).find("[name=keyword]").val() == ''){
            alert('请输入关键词');
        }else{
            searchResult.html('');
            currentPage = 1;
            search()
        }
    });
    function search() {
        var keyword = $("input[name=keyword]").val();
        $.getJSON("/search/" + keyword + "?page=" + currentPage, function(res){
            $.each(res.data, function(i, v){
                searchResult.append('<div class="col-sm-12 col-md-6"><div class="card"><div class="card-header"><a href="' + v.url + '" target="_blank">' + v.title + '</a></div><div class="card-body">' + v.description + '</div><div class="card-footer">' + v.domain + '<a href="search-detail/' + v.id + '" target="_blank">快照</a></div></div></div>');
            });
            timeCost.html(res.cost / 10000000 + 'ms')
        })
    }
});