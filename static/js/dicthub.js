/*!
 * Dicthub JS
 */
function updateDict(event) {
    const form = $('#editWordForm')
    event.preventDefault()
    // 解析數據
    let arr = form.serializeArray()
    console.log(arr)
    let dict = arr2Dict(arr);
    console.log(JSON.stringify(dict))
    // 校驗數據
    if (dict['name'].trim() === '' || dict['desc'].trim() === '') {
        $('#submitTip').css('color', '#dc616a').text("請完成必填項後再提交")
        return
    }
    // 提交數據
    $.ajax({
        url: '/api/dicts/' + dict["id"],
        type: 'PUT',
        data: JSON.stringify(dict),
        dataType: 'json',
        contentType: 'proposal/json; charset=UTF-8',
        complete: function (xhr, resp, text) {
            if (xhr.status === 200) {
                form.trigger('reset');
                const lastPage = $.trim(document.referrer)
                if (lastPage) {
                    history.go(-1)
                }
            } else {
                $('#submitTip').css('color', '#dc616a').text("編輯失敗：" + text)
            }
        }
    });
}

function addWord(event) {
    const form = $('#addWordForm')
    event.preventDefault()
    // 解析數據
    let arr = form.serializeArray()
    // console.log(arr)
    let word = arr2word(arr);
    // 校驗數據
    if (word['writing'].trim() === '') {
        $('#submitTip').css('color', '#dc616a').text("請完成必填項後再提交")
        return
    }
    // console.log(word)
    // 提交數據
    $.ajax({
        url: '/api/words/',
        type: 'POST',
        data: JSON.stringify(word),
        dataType: 'json',
        contentType: 'proposal/json; charset=UTF-8',
        complete: function (xhr, resp, text) {
            // TODO: 不知道爲什麼會報錯誤：parsererror SyntaxError: JSON.parse: unexpected character at line 1 column 1 of the JSON data
            if (xhr.status === 200) {
                submitCount = $('#submitTip').attr('count')
                if (isNaN(submitCount)) {
                    submitCount = 0
                }
                submitCount++
                $('#submitTip').attr('count', submitCount)
                $('#submitTip').css('color', '#17a2b8').text("已成功添加 " + submitCount + " 項")
                // keep catogry value
                catVal = $('#catalogInput select').val()
                // 注意：jquery 中的 form 跟 getElementById 中獲取的對象不同，reset 方法未定義
                form.trigger('reset');
                $('#catalogInput select').val(catVal)
            } else {
                $('#submitTip').css('color', '#dc616a').text("添加失敗：" + text)
            }
        }
    });
}

function updateWord(event) {
    const form = $('#editWordForm')
    event.preventDefault()
    // 解析數據
    let arr = form.serializeArray()
    // console.log(arr)
    let word = arr2word(arr);
    // 校驗數據
    if (word['writing'].trim() === '') {
        $('#submitTip').css('color', '#dc616a').text("請完成必填項後再提交")
        return
    }
    // 提交數據
    $.ajax({
        url: '/api/words/' + word["id"],
        type: 'PUT',
        data: JSON.stringify(word),
        dataType: 'json',
        contentType: 'proposal/json; charset=UTF-8',
        complete: function (xhr, resp, text) {
            // TODO: 不知道爲什麼會報錯誤：parsererror SyntaxError: JSON.parse: unexpected character at line 1 column 1 of the JSON data
            if (xhr.status === 200) {
                // 注意：jquery 中的 form 跟 getElementById 中獲取的對象不同，reset 方法未定義
                form.trigger('reset');
                const lastPage = $.trim(document.referrer)
                if (lastPage) {
                    history.go(-1)
                }
            } else {
                $('#submitTip').css('color', '#dc616a').text("編輯失敗：" + text)
            }
        }
    });
}

function propose(event) {
    const form = $('#proposeForm')
    event.preventDefault()
    // 解析數據
    let arr = form.serializeArray()
    // console.log(arr)
    let proposal = arr2proposal(arr);
    // 校驗數據
    if (proposal['applicant']['email'].trim() === '') {
        $('#submitTip').css('color', '#dc616a').text("請完成必填項後再提交")
        return
    }
    // 提交數據
    $.ajax({
        url: '/api/proposals/',
        type: 'POST',
        data: JSON.stringify(proposal),
        dataType: 'json',
        contentType: 'proposal/json; charset=UTF-8',
        complete: function (xhr, resp, text) {
            // TODO: 不知道爲什麼會報錯誤：parsererror SyntaxError: JSON.parse: unexpected character at line 1 column 1 of the JSON data
            if (xhr.status === 200) {
                $('#staticBackdrop').modal('show')
                // 注意：jquery 中的 form 跟 getElementById 中獲取的對象不同，reset 方法未定義
                $('#modalOkBtn').click(function(){
                    form.trigger('reset');
                    const lastPage = $.trim(document.referrer)
                    if (lastPage) {
                        history.go(-1)
                    }
                })
            } else {
                $('#submitTip').css('color', '#dc616a').text("提交失敗: " + text)
            }
        }
    });
}


function arr2Dict(arr) {
    let dict = {}
    for (let i = 0; i < arr.length;i++) {
        const v = arr[i]
        if (['preferSpecLinkIds', 'tags'].indexOf(v.name) != -1) {
            dict[v.name] = v.value.split(',')
        } else if (['id', 'name', 'desc', 'cover', 'feedbackEmail', 'catalogText', 'specText']) {
            dict[v.name] = v.value.trim();
        }
    }
    return dict
}

function arr2word(arr) {
    let word = {}
    let specs = []
    for (let i = 0; i < arr.length;) {
        const v = arr[i]
        if (v.name === 'tags') {
            word[v.name] = v.value.split(' ').filter(tag => tag.trim() !== '')
            i++
        } else if (v.name === 'linkId') {
            if (v.value !== '' && arr[i + 1].value.trim() !== '') {
                specs.push({
                    linkId: v.value,
                    value: arr[i + 1].value.trim(),
                    note: arr[i + 2].value.trim()
                })
            }
            i += 3
        } else if (v.name === 'catalogLinkId') {
            ids = []
            lastId = ''
            for (id of v.value.split('/')) {
                curId = lastId === '' ? id : lastId + "/" + id
                ids.push(curId)
                lastId = curId
            }
            word['catalogLinkIds'] = ids
            i++
        } else if (['id', 'dictId', 'writing', 'meaning', 'completion'].indexOf(v.name) != -1) {
            word[v.name] = v.value.trim();
            i++
        }else{
            i++
        }
    }
    word['specs'] = specs
    return word
}

function arr2proposal(arr) {
    console.log(arr)
    let proposal = {}
    let applicant = {}
    let kind = ''
    for (let i = 0; i < arr.length;i++) {
        const v = arr[i]
        if(v.name === 'kind'){
            kind = v.value
        }else if (v.name === 'applicantName'){
            applicant['name'] = v.value
        }else if (v.name === 'applicantEmail'){
            applicant['email'] = v.value
        }
    }
    switch (kind) {
        case 'add_word':
        case 'edit_word':
            word = arr2word(arr)
            proposal['word'] = word
            proposal['kind'] = kind 
            break
    }
    proposal['applicant'] = applicant
    return proposal
}

function addSpec() {
    $('#specsInput').append($('#specsInput').children("div").last()[0].cloneNode(true))
    $('#specsInput').find('button').last().removeAttr('disabled')
}
