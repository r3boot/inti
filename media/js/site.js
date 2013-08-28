// BEGIN: LocalStorage functions
function store_get(k) {
    try {
        v = s.get(k)
        console.log('store_get: '+k+' = '+v)
        if (v == 'false') {
            return false
        } else if (v == 'true') {
            return true
        }
        return JSON.parse(v)
    } catch (err) {
        console.log('store_get: could not retrieve '+k)
        return null
    }
}

function store_set(k, v) {
    console.log('store_set: '+k+' = '+JSON.stringify(v))
    s.set(k, JSON.stringify(v))
    s.save()
}

function reset_store() {
    console.log('resetting LocalStorage')
    store_set('spots', {})
    store_set('groups', {})
}


// BEGIN: Application initialization
function load_configuration() {
    $.ajax({
        url: '/config',
        type: 'get',
        success: function(response) {
            store_set('spots', response['spots'])
            load_spots(response['spots'])

            store_set('groups', response['groups'])
            //load_groups(response['groups'])
        },
        error: function(xhr, textStatus, errorThrown) {
            console.log('Failed to load configuration: '+textStatus+', '+errorThrown)
        }
    })
}

function load_spots(spots) {
    var content = '<h4>Available spots</h4>'
    content += '<div class="row">'
    content += '<div class="span1"><b>Name</b></div>'
    content += '<div class="span3"><b>Description</b></div>'
    content += '<div class="span2"><b>Location</b></div>'
    content += '<div class="span1"><b>Color</b></div>'
    content += '</div>'
    for (var spot_name in spots) {
        var spot_info = spots[spot_name]
        content += '<div class="row">'
        content += '<div class="span1">'+spot_name+'</div>'
        content += '<div class="span3">'+spot_info['description']+'</div>'
        content += '<div class="span2">'+spot_info['location']+'</div>'
        content += '<div class="span1">#000000</div>'
        content += '</div>'
    }
    $('#available_spots').html(content)
}

function setup_application() {
    var content = '<div class="tabbable">'
    content += '<ul class="nav nav-tabs">'
    content += '<li class="active"><a href="#overview" data-toggle="tab">Overview</a></li>'
    content += '<li><a href="#mixer" data-toggle="tab">Mixer</a></li>'
    content += '</ul>'
    content += '<div class="tab-content">'
    content += '<div class="tab-pane active" id="overview">'
    content += '</div>'
    content += '<div class="tab-pane" id="mixer">'
    content += '<p>Howdy, I am in Section 2.</p>'
    content += '</div>'
    content += '</div>'
    content += '</div>'
    $('#content').html(content)

    // Initialize all tabs
    content = '<p />'
    content += '<div class="container" id="available_spots">'
    content += '</div>'
    $('#overview').html(content)

    $('#overview a').click(function (e) {
        e.preventdefault();
        $(this).tab('show');
    })

    content = '<p>'
    content += 'mixer'
    content += '</p>'
    $('#mixer').html(content)

    $('#mixer a').click(function (e) {
        e.preventdefault();
        $(this).tab('show');
    })

}

function main() {
    $(document).ready(function() {
        s = new Persist.Store('inti')
        reset_store()
        load_configuration()
        setup_application()
    })
}

/* GO GO GO */
main()
