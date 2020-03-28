## 从orm学习接口

#### 声明

```type Ormer interface {
	//
	ReadAll(md interface{}, cols ...string) error
	//
	Insert(md interface{}) error
	//
	Rollback() error
	//
	Update(md interface{}, cols ...string) error
	//
	CloseClient()
	//
	Delete(md interface{}, cols ...string) error
	//
	QueryTable(ptrStructOrTableName interface{}) (qs *querySet)
}
```
如上是一个接口声明，一般来说接口声明遵循以下几点：
* 接口命名一般以er结尾
* 接口定义的内部方法声明不用func引导
* 接口定义只有声明没有实现

#### 调用

未初始化的接口调用会panic

#### 动态类型和静态类型

* 动态类型
  接口绑定的具体实例的类型称为接口的动态类型
* 静态类型
  接口被定义时确定的类型叫接口的静态类型,静态类型的本质就是接口的方法签名集合。


#### 优点 

* 解耦 
  Go的非侵入式接口使得层与层之间的代码更加干净，具体类型和接口的实现之间不需要显式声明，增加了接口使用的自由度。
* 实现泛型
  使用空接口作为函数或者方法参数能够用在需要泛型的场景中。例如beego/orm中使用了大量的空接口通过反射实现了泛型编程。


#### 泛型

#### 底层实现
