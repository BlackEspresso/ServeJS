namespace js2table {

    export function generate(obj: any) {
        var objType = getType(obj);
        return buildHtml(obj);
    }

    function buildHtml(obj) {
        let type = getType(obj)
        let func = getBuildFunc(type);
        return func(obj);
    }

    function getBuildFunc(type) {
        if (type == 'array') {
            return buildArray;
        } else if (type == 'object') {
            return buildObject;
        } else {
            return buildValue;
        }
    }

    function getAllProperties(obj) {
        var props = [];
        for (var x = 0; x < obj.length; x++) {
            var objProps = getProperties(obj[x]);
            if (x == 0)
                props = objProps;
            else {
                for (var y = 0; y < objProps.length; y++) {
                    if (props.indexOf(objProps[y]) == -1)
                        props.push(objProps[y]);
                }
            }
        }
        return props;
    }

    function extendObject(obj,props) {
        for (var x = 0; x < props.length; x++) {
            if (obj[props[x]] === undefined) {
                obj[props[x]] = toString(null);
            }
        }
    }

    function getProperties(obj) {
        var props = [];
        for (var x in obj) {
            props.push(x);
        }
        return props;
    }

    function buildArray(obj) {
        var str = '<table>';
        //var sameProps = hasSameProperties(obj);
        if (getType(obj) == 'array' && obj.length > 1 && getType(obj[0]) == 'object') {
            var props = getAllProperties(obj);
            str += '<tr>';
            for (var x = 0; x < props.length; x++) {
                str += '<th>' + props[x] + '</th>';
            }
            str += '</tr>';
            for (var x = 0; x < obj.length; x++) {
                extendObject(obj[x], props);
                let val = buildArrayRow(obj[x]);
                str += val;
            }
        } else {
            str += '<tr>';
            for (var x = 0; x < obj.length; x++) {
                let val = buildHtml(obj[x]);
                str += '<td>' + val + '</td>';
            }
            str += '</tr>';
        }
        str += '</table>';
        return str;
    }

    function buildArrayRow(obj) {
        var str = '<tr>';
        for (var x in obj) {
            var val = buildHtml(obj[x]);
            str += '<td>' + val + '</td>';
        }
        str += '</tr>';
        return str;
    }

    function buildObject(obj) {
        var str = '<div>';
        for (var x in obj) {
            let val = buildHtml(obj[x]);
            str += '<strong>' + x + '</strong><div>' + val + '</div>'
        }
        str += '</div>';
        return str;
    }

    function buildValue(obj) {
        return toString(obj);
    }

    function toString(obj): string {
        if (obj == null)
            return '';
        return obj.toString();
    }

    function getType(obj: any) {
        if (Array.isArray(obj))
            return 'array';
        return typeof obj;
    }
}
