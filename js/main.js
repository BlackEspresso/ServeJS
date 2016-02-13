function onStart(){
	var m = require('httpmappings')
	var t = require('tasks')
	var c = require('cache')
	
	m.addMapping('/websocket','websocket');
	m.addMapping('/writefile','writefile');
	t.startTasks();
	
	var yaml = require('settings');
	var conf = yaml.read('./serverjs.yaml');
	c.set('settings', JSON.stringify(conf))

	var server = require('httplistener')
	server.requestFuncName = 'onRequest';
	//server.addr = ':8080'
	server.startAndWait()
}

function onRequest(resp,req){
	resp.contentType='text/html'
	//return JSON.stringify(req.cookies);
	router(resp,req)
		.on('/gquery',gquery)
		.on('/hello',hello)
		.on('/showcache',showcache)
		.on('/mailto', mailTo)
		.on('/run', run)
		.on('/login',login)
		.on('/register',register)
		.on('/req',request)
		.on('/editjs',editjs)
		.on('/boerse',boerse)
		.on('/scandns',scandns)
		.on('/schedule',schedule)
		.on('/datetime',datetime)
		.on('/cache',cacheFunc)
		.on('/header',header)
		.on('/htmlcheck',htmlCheck)
		.on('/mongo',mongo)
}

function showcache(resp,req){
	var cache = require('cache');
	if(!isAuth(req)){
		redirect(resp,'/login')
		return
	}
	resp.write(JSON.stringify(cache.all()))
}

function redirect(resp,url){
	resp.statusCode=302;
	resp.header = {
		'Location':url
	}
}

function mongo(resp,req){
	var m = require('mongodb')
	var sess = m.newSession('localhost')
	var c= sess.DB('jstest').C('test')

	c.RemoveAll({})

	c.Insert({'Name':'Tester1','Street':'here'})

	c = c.Find({'Name':'Tester1'})
	for(var x in c){
		resp.write(x + '\n');
	}

	var test = m.all(c.Select({'Name':1}));

	resp.write(JSON.stringify(test));
}

function gquery(resp,req){
	var goquery = require('goquery')
	var http = require('http')
	var cResp = http.do({url:'http://google.com'});
	var doc = goquery.newDocument(cResp.val.body)
	var links = doc.ExtractAttributes('body');
	resp.write(JSON.stringify(links))
}

function htmlCheck(resp,req){
    var url = req.formValues.url;
    if(url==null){
        resp.write('url is empty')
        return
    }
    var cResp = http.do({url:url});
    htmlcheck.loadTags('./static/tags.json')
    
	var err = htmlcheck.validate(cResp.val.body)
	resp.write(JSON.stringify(err))
}

function header(resp,req){
	resp.header = {
		'Set-Cookie':'c=5',
		'Test-Header':'4'
	}
	resp.write('4')
}

function cacheFunc(resp,req){
	var key = req.formValues.key;
	var val = req.formValues.value;
	if (val == null){
		resp.write(cache.get(key))
	}else{
		cache.set(key,val);
		resp.write('ok')
	}
	cache.save('./static/cache.json')
} 

function datetime(resp,req){
	resp.write(new Date().toString())
}

function boerse(resp,req){
    var symbol = req.formValues.symbol;
    var url = 'http://finance.yahoo.com/webservice/v1/symbols/'+symbol+'/quote?format=json';
    var siteResp = http.do({url:url});
    //resp.write(JSON.stringify(siteResp))
    var serviceResp = JSON.parse(siteResp.val.body);
    //resp.write(JSON.stringify(JSON.parse(siteResp.body),null,2))
    resp.write(serviceResp.list.resources[0].resource.fields.name);
    resp.write('\n')
    resp.write(serviceResp.list.resources[0].resource.fields.price);
}

function scandns(resp,req){
    var fileRead = readFile('static', 'subdomain_wordlist.txt');
    var lines = fileRead.Suc.split('\n');

    if (req.formValues.host === null) {
        resp.write('error')
    }

    var count = 5;
    if (req.formValues.count !== null) {
        count = parseInt(req.formValues.count)
    }

    var ret = {};
    var regExIp = /([^\s]+)\s+(\d+)\s+(\w+)\s+(\w+)\s+(\d+\.\d+\.\d+\.\d+)/ig;

    resp.write('[\n')
    for (var x = 0; x < lines.length&&x<count; x++) {
        var dnsReponse = resolve(lines[x].trim() + '.' + req.formValues.host);
        var entries= dnsReponse.split('\n');
        
        for(var z=0;z<entries.length;z++){
    		var matches= entries[z].split(/\s+/);
        	for (var y = 0; y < matches.length; y++) {
            	if (y > 4)
                	matches[4] += matches[y];
            }
        	matches.splice(5, matches.length - 5);
        	if(matches.length === 1 && matches[0]==="")
        		continue;
            resp.write(JSON.stringify(matches));
            if(x+1<lines.length && x+1<count){
                resp.write(',\n');
            }
        }
    }
    resp.write('\n]');
}

function editjs(resp,req){
	var templ = require('templating');
	if(!isAuth(req)){
		redirect(resp,'/login');
		return;
	}
	var file = require('file')
	templ.reloadTemplates();
	var js = file.readFile('./js/main.js');
	var tmpl = templ.runTemplate("EditJs.thtml",{MainJs:js.val});
	resp.write(tmpl)
}

function onWebSocketRequest(message){
	console.log('js: '+message)
}

function request(resp,req){
	siteResp = httpDo({url:'http://google.com',method:'GET'})
	resp.write(JSON.stringify(siteResp))
}

function mailTo(resp,req){
	var receiver = req.formValues.mailto
	var k = send(receiver,'hallo',JSON.stringify(req))
	resp.write(k.Err)
}

function run(resp,req){
	resp.contentType='text/plain'
	var cmd = runCmd('echo','4','5')
	if(cmd.val)
		resp.write(cmd.val);
	else
		resp.write(cmd.error);
}

function hello(resp,req){
	resp.write('hello')
}

function schedule(resp,req){
	var d = new Date();
	d.setMinutes(d.getMinutes()+1);

    addTask('test',d,function(){
        //var cmd = runCmd('cmd','/c','dir');
        console.log('task started: ' + new Date())
    });
	
	resp.write('scheduled')
}

function getCollection(mongodb){
	var sess = mongodb.newSession('localhost')
	var c = sess.DB('jstest').C('test')
	return c;
}

function register(resp,req){
	var mongodb = require('mongodb')
	var c = getCollection(mongodb);

	var templating = require('templating');
	templating.reloadTemplates();

	if(req.method=='GET'){
		var templ = templating.runTemplate("login.thtml",{ActionName:'register'});
		resp.write(templ)
	}
	if(req.method=='POST'){
		var username = req.formValues.username[0];
		var password = req.formValues.password[0];
		c.Insert({'Username':username,'Password':password});
	}
}


function getCookie(req,name){
	var cookies = {};
	if(req.cookies!=null){
		for(var x = 0;x<req.cookies.length;x++){
			if(req.cookies[x].Name==name)
				return req.cookies[x];
		}
	}
	return undefined;
}

function getAuth(req){
	var cache = require('cache');
	var cookie = getCookie(req,'userid')
	if (cookie == null)
		return false;
	var userId = cookie.Value;
	var fromCache = cache.get('userid_'+userId);
	if(fromCache == null)
		return false;
	var obj = JSON.parse(fromCache);
	return obj;
}

function isAuth(req){
	var auth = getAuth(req);
	return auth!==null && auth.isAuth;
}

function login(resp,req){
	var templating = require('templating');
	templating.reloadTemplates();

	if(req.method=='GET'){
		var auth = getAuth(req);
		if(auth != null && auth.isAuth){
			message = 'hello ' + auth.username
		} else{
			message = 'please log in'
		}

		var templ = templating.runTemplate("login.thtml",
			{ActionName:'login',Message:message});
		resp.write(templ)
	}
	else{
		var mongodb = require('mongodb')
		var cache = require('cache')
		var c = getCollection(mongodb);

		var username = req.formValues.username[0];
		var password = req.formValues.password[0];

		c =c.Find({'Username':username})
		var user = mongodb.one(c);
		if(user !== false && user.Password==password){
			var auth = getAuth(req);
			if(auth!=null && auth.id!=null){
				cache.remove('userid_'+auth.id)
			}
			var id = (Math.random()*100000000)|1;

			cache.set('userid_'+id, JSON.stringify({
				isAuth:true,
				id:id,
				date: new Date(),
				isAdmin:user.isAdmin === true,
				username:username})
			);
			resp.header = {'Set-Cookie':'userid='+id+';HttpOnly'};
		}
		else
			resp.write('not ok');
	}
}

function router(resp,req){
	var obj= {
		on:function(url,func){
			if(req.url.indexOf(url)===0){func(resp,req)}
			return obj;
		}
	};
	return obj;
}






