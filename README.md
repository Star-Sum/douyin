# 极简版抖音开发

## 项目介绍
基于gin-gorm架构的极简抖音后端项目  

项目地址：[https://github.com/Star-Sum/douyin](https://github.com/Star-Sum/douyin)

## 项目分工
| 团队成员 | 主要贡献                       |
|------|----------------------------|
| 任磊   | 项目整体架构设计、技术选型、消息发送功能实现     |
| 王宇杰  | Feed流、投稿功能实现、部分辅助组件实现      |
| 王巍良  | 用户关注及关注列表功能实现，smartart分支   |
| 李承泽  | 粉丝以及好友列表功能实现、关注的master分支实现 |
| 刘志澎  | 用户注册、用户登录、用户信息接口实现         |
| 王轶则  | 查询聊天记录、发送消息接口实现            |
| 齐会英  | 根据数据库接口实现点赞功能和点赞列表拉取       |
| 魏贻辉  | 利用数据库接口实现评论功能及获取视频下评论列表    |

## 项目实现
### 技术选型与相关开发文档
- Gin
- Gorm
- FFmpeg
- Mysql 
- Redis

### 架构设计
本项目采用MVC分层设计模型分离了数据访问层、实体层、业务层、控制层，
从而降低了代码的耦合度，提高了项目的可维护性。

本项目采用redis作为缓存层、mysql作为持久层、Gin作为Web框架。

### 项目代码介绍
#### 视频Feed流实现
缓存层Redis中,
保存用户Feed流收件箱，在每次视频发布时都会将视频插入关注者的sortedList中

***Feed流推送功能***

Feed流的模式分为拉模式、推模式和推拉结合模式。拉模式即每个粉丝从被关注者的“信箱”中主动拉取视频信息。这种Feed流模式适用于有大量粉丝群体的用户；推模式即每个用户在发布视频后都向其粉丝的“收件箱”中投递。这种Feed流适用于存在粘性用户的情况。而推拉结合即是将上述两种情况结合起来，在不同情况下使用不同的Feed流模式。

本项目使用推模式，即在每个用户发布视频时即向其粉丝的收件箱内都分别推送视频信息。具体实现上使用了Redis的sortedList功能，以时间戳为排序socre，实现倒叙推送视频流，在未登录情况下按照时间发布倒叙依次推送视频信息，在登录状态下推送该用户关注者视频。为了保证视频流信息的完整性与及时性，每次登录都会初始化登录用户的Feed流收件箱，拉取当前数据库内的最新30条关联视频信息。

在拉取视频流时首先对token进行鉴权（如果处于登录状态），如果鉴权成功则根据请求用户的UID和请求发出的时间戳拉取最近的30条视频信息，在查询时采用旁路缓存策略。先查询Redis 缓存，如果未能命中，则查询数据库，后序将信息更新到缓存中。

#### 用户注册、用户登录、用户信息接口实现
***UserInfo模型***  

持久化层Mysql中，
包含：ID、Avatar、BackgroundImage、FavoriteCount、IsFollow、Name、Signature、TotalFavorited、WorkCount, 其中ID作为主键。  

缓存层Redis中，
记录用户的信息，把UserInfo结构体中的每个字段存储到hash中

***UserAccountInfo模型***  

持久化层Mysql中，
包含：UserId、Password、Username，其中UserId为主键

***用户注册功能***

在注册之前先向数据库确认Username不存在，确认Username唯一后分别向UserInfo表和UserAccountInfo表添加记录，ID使用雪花算法升成，Password则用MD5加密。注册成功后，生成包含UserID的令牌返回给客户端。

***用户登录功能***

从请求中解析出Username和Password参数，把Password用MD5加密，从UserAccountInfo表中查询是否存在相同的记录。如果不存在，登录失败；如果存在，则登录成功，并生成包含UserID的令牌返回给客户端。

***用户信息功能***

从请求中解析出Token和UserID，首先对用户鉴权，Token解析错误立即返回，成功后则进行下一步操作。

根据UserID从Redis中查询是否存在用户信息。若缓存命中，则从Redis中获取信息直接返回；若缓存未命中，则从MySQL中获取信息，返回用户信息并把用户信息缓存到Redis中。

#### 视频投稿功能
***PublishInfo模型***  

持久化Mysql中，
包含：UrlRoot、UseID、Time、Title、VedioID、OpID，其中VedioID作为自增主键。

缓存层Redis中，
记录投稿后的视频信息（key:"Publish:_"+UID,value:视频信息的序列化json字符串）

***发布视频功能***

在发布视频时先进行token鉴权，判断登录状态是否有效。token鉴权完成后即读取传入文件基本信息，因为数据库接口要求必须直接给出视频url，所以本处直接采取静态资源映射的方式。由于视频发布需要提供封面，而前端并未提供封面上传功能，此处使用ffmpeg的截图功能，默认截取每个视频的第一帧作为封面，同样采取静态资源映射的方式。在url信息转换成功后写入数据库与缓存，同时更新用户的作品数。

***拉取视频列表功能***

由于我们需要在个人页获取每个用户的投稿信息，就需要我们根据每个用户的UID来拉取发布视频的列表。我们首先从缓存中拉取视频流信息，当缓存中信息量不足时我们采用旁路存取策略从持久化存储中拉取投稿视频信息并将其写入缓存中。

#### 点赞以及喜欢列表

***点赞操作***

首先对用户鉴权，并从Token中解析出当前用户fromUserID 。传入的ActionType参数是string，先解析为int64——如果actionType不是1（点赞）或者2（取消点赞），然后判断数据库中是否存在LikeVideoID的视频如果存在则返回结构体中状态码。

***喜欢列表***

首先判断token是够合法，提取fromUserID,如果ID视频状态码为1获取点赞过的视频。

#### 评论及评论列表接口实现
***Comment模型***  

持久化层Mysql中，  
包含：ID、VedioId、UsrID、Content、CreateData，其中ID作为主键。

缓存层Redis中   
记录视频下的评论列表（"Comment:videoId_" + strconv.FormatInt(videoId, 10)，评论列表的json字符串）

***评论功能***

1.根据token进行登录校验

2.传入的ActionType参数是string，先判断是否合法（是否为1或2）
* 如果ActionType是1，对应为添加评论操作，则根据雪花算法生成评论iD，并将Comment结构体写入持久层。
* 如果ActionType是2，对应为删除评论操作，验证评论ID是否存在后执行删除操作，并写入持久层以及缓存层。

#### 关注及关注列表实现
***Relation模型***
持久化层Mysql中，  
包含：ID、UserID、ToUserID、Time，其中ID作为自增主键。

缓存层Redis中，  
记录用户的关注信息（key : "Relation:_follow" + UserID + "to" + ToUserID, value：关注信息的序列化json字符串）  
用户的关注列表信息（key : "Relation:_focus" + UserID , value：关注列表信息的序列化json字符串）

***关注功能***

在进行关注操作前，首先对用户鉴权，并从Token中解析出当前用户ID。
判断ActionType是否合法，根据ActionType进行关注或取消关注的操作。

- 关注时：基于gorm进行Mysql操作，新建关系，对数据库进行更新，后序会将数据同步到Redis缓存中。取消关注操作与关注操作类似。
- 查询关注列表时：采用旁路缓存策略。先查询Redis 缓存，如果未能命中，则查询数据库，后序将信息更新到缓存中。为了保证数据一致性，当有数据变化时，先更新数据库，随后删除旧的缓存。后序查询未命中后会更新缓存。

#### 粉丝列表以及好友列表的实现
***Relation模型***

持久化层Mysql中，  
包含：ID、UserID、ToUserID、Time，其中ID作为自增主键。

缓存层Redis中，  
记录用户粉丝列表信息（key : "Relation:_fans" + UserID value：粉丝列表信息的序列化json字符串）  
用户的好友列表信息（key : "Relation:_follow" + UserID +"to"+ToUserId, value：好友列表信息的序列化json字符串）

***粉丝列表功能***

查询粉丝列表的时候首先对用户进行鉴权，然后会从缓存中首先查找是否包含key的相关内容，如果有就直接返回，如果没有就从数据库中查询to_user_id相关字段对应的user_id字段，然后通过user_id查询用户信息表，将用户信息封装形成列表返回，并将其放到redis中进行存储，方便之后的查询。

***好友列表功能***
查询好友列表时首先对用户进行鉴权操作，从缓存中查找是否有相关信息，如果缓存命中就返回，没有命中就从数据库中查询当user_id中存在to_user_id并且to_user_id中存在user_id字段（互相关注）信息，通过to_user_id字段查询用户信息表，并将用户信息封装形成列表返回，并加载到redis中存储，方便之后的查询。

***关注操作***

当用户进行关注或者取关操作的时候首先对用户进行鉴权，然后通过actiontype字段（1:关注，2:取关）判断对用户关注还是取关。

运行之后通过管道操作先将正确信息返回给用户（保证用户的体验）并将信息写入redis中进行存储，然后设置时间每隔一段时间将redis中的关注信息持久化到数据库中并将redis中有关关注者以及被关注者的redis全部删除，下一次获取从数据库中获取后写入redis保证数据一致性。

***关注列表***

根据key先查询Redis 缓存，如果未能命中，则通过user_id查询数据库，后序将信息更新到缓存中。

#### 查询聊天记录以及发送消息的实现

***Message模型***

持久化Mysql中，  
包含：Content、CreateTime、FromUseID、ID、ToUserID，其中ID作为自增主键。

缓存层Redis中，  
发送的聊天记录键值对格式（key:"messages:user_"+toUserID+"_from_"+fromUserID,value:聊天信息的序列化json字符串）

***查询聊天记录功能***

1.在查询聊天记录时先进行token鉴权，判断登录状态是否有效，并token中读到发送者的ID fromUserID；

2.传入的ToUserID参数类型是string，先解析为int64。

3.判断fromUserID和toUserID是否是同一个。

4.验证数据库中是否存在ID为toUserID的用户。

5.先获取自己发给对方的聊天记录列表，再获取对方发给自己的聊天记录列表，然后将二者合并并按时间戳排序，得到最终的两人之间的聊天记录列表。每次查询中，都先从缓存中查询聊天记录列表，当缓存未命中时用旁路存取策略从持久化存储中查询聊天记录列表，并将其序列化写入缓存中。

***发送信息功能***

1.在查询聊天记录时先进行token鉴权，判断登录状态是否有效，并token中读到发送者的ID fromUserID。

2.传入的ActionType参数是string，先解析为int64

3.传入的ToUserID参数类型是string，先解析为int64。

4.判断fromUserID和toUserID是否是同一个。

5.验证数据库中是否存在ID为toUserID的用户。

6.获取消息内容和当前时间戳，将消息插入数据库，判断数据库返回结果是否正确。

## 测试结果

<img src="./Markdown/show1.jpg" >
<img src="./Markdown/show2.jpg" >
<img src="./Markdown/show3.jpg" >
<img src="./Markdown/show4.jpg" >


## Demo 演示视频
[show.mp4](https://github.com/Star-Sum/douyin/blob/master/show.mp4)

## 项目总结与反思

1.本次在项目开发过程中使用了MVC架构，极大地降低了项目耦合程度，为项目开发带来方便

2.本次采用redis缓存层+mysql持久层的开发思路，加快了数据的存取，提高了效率

3.本次架构尚属于单体架构的范畴，预计在未来实现从单体架构到微服务架构的转变

4.对于一些极端情况下的数据未进行测试

5.本次在使用雪花算法实现UID生成时没有考虑数据中心ID和机器ID的生成，导致有极小可能存在UID生成重复的现象。这里我们预计在未来通过Redis+分布式中心的方式实现上述设想

