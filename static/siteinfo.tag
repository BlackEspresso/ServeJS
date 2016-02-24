<siteinfo>
  <h3>Status: {siteinfo.status}</h3>
  <h3>Header</h3>
  <jobject data={siteinfo.header}></jobject>
  <h3>Invalid Tags</h3>
  <jtable data={siteinfo.invalidTags} props={invalidTagsProps}></jtable>
  <h3>Cookies</h3>
  <jtable data={siteinfo.cookies} props={cookieProps}></jtable>
  <h3>Hyperlinks</h3>
  <jtable data={siteinfo.allUrls} clickHandler={scanSite} props={hrefProps}></jtable>
  <script>
    this.siteinfo = opts.siteinfo || {}
    this.cookieProps=[{name:'Name'},{name:'Value'}]
    this.invalidTagsProps = [{name:'TagName'},
    {name:'AttributeName'},{name:'ReasonText'}]
    this.headerProps = [{name:'Name'}]
    this.hrefProps = [{name:"href",type:'link'}]
    if (this.siteinfo.invalidTags){
      invalidTags(this.siteinfo.invalidTags);
    }

    tableClick(e){
      var el = e.target;
      var href = el.attributes.getNamedItem('href');
      if(href!=null){
        e.preventDefault();
        $('#inputUrl').val(href.value);
        $('#formScanUrl').submit();
      }
    }

    function invalidTags(tags){
      for(var x=0;x<tags.length;x++){
        var text = '';
        switch(tags[x].Reason){
          case 0:
            text='Invalid tag';
            break;
          case 1:
            text='Invalid attribute';
            break;
          case 2:
            text = 'Closed before opened';
            break;
          case 3:
            text = 'Not properly closed';
            break;
          case 4:
            text = 'duplicate attribute';
        }
        tags[x].ReasonText = text;
      }
    }

    </script>
</siteinfo>

<jobject>
  <table>
    <tr each={k,v in data}>
      <td>{k}</td><td>{v}</td>
    </tr>
  </table>
  <script>
    this.data = opts.data||{};

    this.on('update',function(){
      //this.rows.length=0;
      //for(var x=0;x<this.data.length;x++){

      //}
    })

  </script>
</jobject>

<jtable>
  <table onclick={parent.tableClick}>
    <tr>
      <th each="{m in props}">{m.displayName||m.name}</th>
    </th>
    <tr each={row in rows}>
      <td each={k in row.cells}>
        <a href="{k.val}" if="{k.type=='link'}">{k.val}</a>
        <span if="{k.type==null}">{k.val}</span>
      </td>
    </tr>
  </table>

  <script>
    this.props= opts.props||{};
    this.data = opts.data||{};
    this.rows = [];

    this.on('update',function(){
      this.rows.length=0;
      for(var x=0;x<this.data.length;x++){
        var row = {cells:[]};
        for(var p=0;p<this.props.length;p++){
          row.cells.push({val:this.data[x][this.props[p].name]||null,
            type:this.props[p].type,
            class:this.props[p].class
          });
        }
        this.rows.push(row);
      }
    })

  </script>
</jtable>
