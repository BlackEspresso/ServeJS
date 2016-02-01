function onRequest(resp,req){
	resp.contentType='text/html'
	//return JSON.stringify(req.cookies);
	router(resp,req)
		.on('/login',login)
		.on('/admin',admin)
}

function router(resp,req){
	var obj= {
		on:function(url,func){
			if(url==req.url){func(resp,req)}
			return obj;
		}
	}
	return obj;
}

function login(resp,req){
	//reloadTemplates();
	if(req.method=='GET'){
		//resp.write(template("login.thtml","test"))
		//resp.write("hello")
		resp.write(readDir('static'))
	}
	else{
		//resp.write('logged in. no. kidding. password wrong')
	}
}

function admin(resp){
	resp.write("todo")
}
